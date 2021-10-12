package storage_test

import (
	. "AppMetadataAPIServerGo/model"
	. "AppMetadataAPIServerGo/storage"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCobraDB(t *testing.T) {
	var cobraDB = CobraDB{}
	var mockCobraSearch = new(MockSearchEngine)
	cobraDB.Init(mockCobraSearch)

	//test Create and GetBulk
	inputs := getTestInputs()
	for _, input := range inputs {
		mockCobraSearch.On("IndexMetadata", input)
		cobraDB.Create(input)
	}

	actual := cobraDB.GetBulk([]AppMetadataKey{
		{inputs[0].Title, inputs[0].Version},
		{inputs[1].Title, inputs[1].Version},
	})

	assert.Equal(t, actual, inputs)
	mockCobraSearch.AssertExpectations(t)

	//test Query
	var queryParam = QueryParameter{Title: "app1"}
	var mockQueriedKeys = []AppMetadataKey{
		{inputs[0].Title, inputs[0].Version},
	}
	mockCobraSearch.On("QueryMetadata", queryParam).Return(mockQueriedKeys)

	actual = cobraDB.Query(queryParam)
	assert.Equal(t, actual, []AppMetadata{inputs[0]})
	mockCobraSearch.AssertExpectations(t)
}

func TestCobraDB_Create_Duplicate(t *testing.T) {
	var cobraDB = CobraDB{}
	var mockCobraSearch = new(MockSearchEngine)
	cobraDB.Init(mockCobraSearch)
	input := getTestInputs()[0]
	mockCobraSearch.On("IndexMetadata", input)
	cobraDB.Create(input)
	err := cobraDB.Create(input)

	assert.Equal(t, err, errors.New("duplicate record exists"))
	mockCobraSearch.AssertExpectations(t)
}

func TestCobraDB_Concurrency(t *testing.T) {
	var cobraDB = CobraDB{}
	var mockCobraSearch = new(MockSearchEngine)
	cobraDB.Init(mockCobraSearch)
	mockCobraSearch.On("IndexMetadata", mock.Anything)

	i:=0
	var failureCount int64
	var successCount int64
	var waitGroupForCreate sync.WaitGroup
	var waitForAllCreate sync.WaitGroup
	var successRecordComapy string
	waitGroupForCreate.Add(1)
	var concurrency = 1000
	for i < concurrency{
		waitForAllCreate.Add(1)
		go func(){
			testData := getTestInputs()[0]
			testData.Company=fmt.Sprint(rand.Int())
			defer waitForAllCreate.Done()
			waitGroupForCreate.Wait()
			err:= cobraDB.Create(testData)
			if err!=nil{
				atomic.AddInt64(&failureCount, 1)
			}else{
				atomic.AddInt64(&successCount,1 )
				successRecordComapy = testData.Company
			}
		}()
		i+=1
	}
	//release for all create goroutine to mimic concurrent requests
	waitGroupForCreate.Done()
	waitForAllCreate.Wait()
	actual := cobraDB.GetBulk([]AppMetadataKey{{getTestInputs()[0].Title,getTestInputs()[0].Version}})
	assert.Equal(t, actual[0].Company, successRecordComapy)
	assert.Equal(t, mockCobraSearch.IndexMetadataArgs, actual[0])
	assert.Equal(t, int(successCount), 1)
	assert.Equal(t, int(failureCount), concurrency-1)

	mockCobraSearch.AssertExpectations(t)
}

type MockSearchEngine struct {
	mock.Mock
	IndexMetadataArgs AppMetadata
}

func (m *MockSearchEngine) IndexMetadata(metadata AppMetadata) {
	m.IndexMetadataArgs = metadata
	m.Called(metadata)
}

func (m *MockSearchEngine) QueryMetadata(queryParams QueryParameter) []AppMetadataKey {
	args := m.Called(queryParams)
	return args.Get(0).([]AppMetadataKey)
}

func (m *MockSearchEngine) Init() {
	m.Called()
}
