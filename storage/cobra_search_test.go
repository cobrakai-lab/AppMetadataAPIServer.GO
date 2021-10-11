package storage

import (
	. "AppMetadataAPIServerGo/model"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCobraSearch_SimpleCase(t *testing.T) {
	var cobraSearch = cobraSearch{}
	cobraSearch.initInvertedIndex()

	input:= getTestInputs()[0]
	cobraSearch.IndexMetadata(input)
	actual := cobraSearch.QueryMetadata(QueryParameter{License: "mit"})

	assert.Equal(t, len(actual), 1)
	assert.Equal(t, actual[0], AppMetadataKey{input.Title, input.Version})
}

func TestCobraSearch_EmptyContent(t *testing.T) {
	var cobraSearch = cobraSearch{}
	cobraSearch.initInvertedIndex()
	actual := cobraSearch.QueryMetadata(QueryParameter{License: "mit"})
	assert.Equal(t, len(actual), 0)
}

func TestCobraSearch_MultipleMatch(t *testing.T){

	var cobraSearch = cobraSearch{}
	cobraSearch.initInvertedIndex()
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
	assert.Equal(t, actual, expected)

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
