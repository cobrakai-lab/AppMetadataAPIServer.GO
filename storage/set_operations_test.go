package storage_test

import (
	. "AppMetadataAPIServerGo/storage"
	"reflect"
	"testing"
)

func TestIntersectAll(t *testing.T) {
	testDataCollections:= createTestDataCollection()
	for _, testData := range testDataCollections{
		input := testData.input
		actual := IntersectAll(input)
		validate(t, input, actual, testData.expected )
	}
}

func validate(t *testing.T, input [][]AppMetadataKey, actual []AppMetadataKey, expected []AppMetadataKey) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual not equal expected. Input: %+v, Actual: %+v, Expected: %+v", input, actual, expected)
	}
}

func createTestDataCollection() []testData{
	return []testData{
		{
			input: [][]AppMetadataKey{},
			expected : []AppMetadataKey{},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app1", Version: "2.0"}},
			},
			expected : []AppMetadataKey{
				{Title: "app1", Version: "1.0"},
				{Title: "app1", Version: "2.0"},
			},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app1", Version: "2.0"}},
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "2.0"}},
			},
			expected : []AppMetadataKey{
				{Title: "app1", Version: "1.0"},
			},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app1", Version: "1.0"}},
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "2.0"}},
			},
			expected : []AppMetadataKey{
				{Title: "app1", Version: "1.0"},
			},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "1.0"}},
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "2.0"},{Title: "app2", Version: "1.0"}},
			},
			expected : []AppMetadataKey{
				{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "1.0"},
			},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "1.0"}},
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "2.0"}},
				{{Title: "app3", Version: "1.0"}, {Title: "app1", Version: "1.0"}},
				{{Title: "app1", Version: "1.0"}},
			},
			expected : []AppMetadataKey{
				{Title: "app1", Version: "1.0"},
			},
		},
		{
			input: [][]AppMetadataKey{
				{{Title: "app1", Version: "1.0"}, {Title: "app2", Version: "1.0"}},
				{{Title: "app1", Version: "1.1"}, {Title: "app2", Version: "2.0"}},
				{{Title: "app3", Version: "1.0"}, {Title: "app4", Version: "1.0"}},
				{{Title: "app1", Version: "2.0"}},
			},
			expected : []AppMetadataKey{},
		},
	}
}

type testData struct{
	input [][]AppMetadataKey
	expected []AppMetadataKey
}