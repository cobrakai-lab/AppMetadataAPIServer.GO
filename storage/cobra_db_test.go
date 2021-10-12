package storage_test

import (
	. "AppMetadataAPIServerGo/model"
	. "AppMetadataAPIServerGo/storage"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
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

type MockSearchEngine struct {
	mock.Mock
}

func (m *MockSearchEngine) IndexMetadata(metadata AppMetadata) {
	m.Called(metadata)
}

func (m *MockSearchEngine) QueryMetadata(queryParams QueryParameter) []AppMetadataKey {
	args := m.Called(queryParams)
	return args.Get(0).([]AppMetadataKey)
}

func (m *MockSearchEngine) Init() {
	m.Called()
}
