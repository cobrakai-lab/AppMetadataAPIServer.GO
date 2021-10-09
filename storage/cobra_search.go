package storage

import (
	. "AppMetadataAPIServerGo/model"
	"log"
)

type SearchEngine interface {
	IndexMetadata(appMetadata AppMetadata)
	QueryMetadata(queryParam QueryParameter) []AppMetadataKey
}

type cobraSearch struct{
	invertedIndex map[string]map[string][]AppMetadataKey
}

func (cobraSearch *cobraSearch) IndexMetadata(appMetadata AppMetadata){
	indexProperty(cobraSearch, appMetadata, "title", appMetadata.Title)
	indexProperty(cobraSearch, appMetadata, "version", appMetadata.Version)
	for _, maintainer:= range appMetadata.Maintainers{
		indexProperty(cobraSearch, appMetadata, "maintainerName", maintainer.Name)
		indexProperty(cobraSearch, appMetadata, "maintainerEmail", maintainer.Email)
	}
	indexProperty(cobraSearch, appMetadata, "company", appMetadata.Company)
	indexProperty(cobraSearch, appMetadata, "website", appMetadata.Website)
	indexProperty(cobraSearch, appMetadata, "source", appMetadata.Source)
	indexProperty(cobraSearch, appMetadata, "license", appMetadata.License)
}

func (cobraSearch *cobraSearch) QueryMetadata(queryParam QueryParameter) []AppMetadataKey{
	//todo
	return []AppMetadataKey{}
}

func (cobraSearch *cobraSearch) initInvertedIndex(){
	if cobraSearch.invertedIndex==nil{
		cobraSearch.invertedIndex = make(map[string]map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["title"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["version"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["maintainerName"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["maintainerEmail"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["company"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["website"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["source"] = make(map[string][]AppMetadataKey)
		cobraSearch.invertedIndex["license"] = make(map[string][]AppMetadataKey)
		log.Println("CobraSearch Inverted Index is initialized")
	}else{
		log.Println("CobraSearch Inverted index is already initialized.")
	}
}

func indexProperty(cobraSearch *cobraSearch, appMetadata AppMetadata, propertyName string , propertyValue string ){
	log.Printf("Indexing %s = %s for app metadata title: %s, version: %s\n", propertyName, propertyValue, appMetadata.Title, appMetadata.Version)
	var propertyIndex = cobraSearch.invertedIndex[propertyName]
	var key = AppMetadataKey{appMetadata.Title, appMetadata.Version}
	subIndex,found := propertyIndex[propertyValue]
	if found{
		subIndex = append(subIndex, key)
	}else{
		propertyIndex[propertyValue] = []AppMetadataKey{key}
	}
}

type QueryParameter struct {
	Title           string
	Version         string
	MaintainerName  string
	MaintainerEmail string
	Company         string
	Website         string
	Source          string
	License         string
}