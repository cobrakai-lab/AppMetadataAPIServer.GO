package storage_test

import (
	. "AppMetadataAPIServerGo/model"
	. "AppMetadataAPIServerGo/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCobraSearch_SimpleCase(t *testing.T) {
	var cobraSearch = CobraSearch{}
	cobraSearch.Init()

	input:= getTestInputs()[0]
	cobraSearch.IndexMetadata(input)
	actual := cobraSearch.QueryMetadata(QueryParameter{License: "mit"})

	assert.Equal(t, len(actual), 1)
	assert.Equal(t, actual[0], AppMetadataKey{input.Title, input.Version})
}

func TestCobraSearch_EmptyContent(t *testing.T) {
	var cobraSearch = CobraSearch{}
	cobraSearch.Init()
	actual := cobraSearch.QueryMetadata(QueryParameter{License: "mit"})
	assert.Equal(t, len(actual), 0)
}

func TestCobraSearch_MultipleMatch(t *testing.T){

	var cobraSearch = CobraSearch{}
	cobraSearch.Init()
	testInputs := getTestInputs()
	for _, input:= range testInputs{
		cobraSearch.IndexMetadata(input)
	}

	actual := cobraSearch.QueryMetadata(QueryParameter{Company: "cobrakai"})
	assert.Equal(t, len(actual), 2)
	expected:= []AppMetadataKey{
		{testInputs[0].Title, testInputs[0].Version},
		{testInputs[1].Title, testInputs[1].Version},
	}
	assert.ElementsMatch(t, actual, expected)

	actual = cobraSearch.QueryMetadata(QueryParameter{Company: "cobrakai", License: "apache"})
	expected = 	[]AppMetadataKey{
		{testInputs[1].Title, testInputs[1].Version},
	}
	assert.Equal(t,len(actual), 1)
	assert.Equal(t, actual, expected)

	actual = cobraSearch.QueryMetadata(QueryParameter{Company: "cobrakai", License: "bsd"})
	assert.Equal(t,len(actual), 0)
}

func getTestInputs() []AppMetadata {
	inputs := []AppMetadata{
		{
			Title:   "app1",
			Version: "1.0",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
			},
			Company:     "cobrakai",
			Website:     "www.overthere.com",
			Source:      "https://github.com",
			License:     " mit ",
			Description: "Nice!",
		},
		{
			Title:   "app2",
			Version: "1.0",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
			},
			Company:     "cobrakai",
			Website:     "www.overthere.com",
			Source:      "https://github.com",
			License:     "apache",
			Description: "Awesome!",
		},
	}

	return inputs
}
