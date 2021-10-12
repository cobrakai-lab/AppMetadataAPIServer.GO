package storage

import (
	. "AppMetadataAPIServerGo/model"
	"log"
	"strings"
)

type SearchEngine interface {
	IndexMetadata(appMetadata AppMetadata)
	QueryMetadata(queryParam QueryParameter) []AppMetadataKey
	Init()
}

type CobraSearch struct{
	// map from propertyName : propertyValue : set of keys
	invertedIndex map[string]map[string]AppMetadataKeySet
}

func (cobraSearch *CobraSearch) IndexMetadata(appMetadata AppMetadata){
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

func (cobraSearch *CobraSearch) QueryMetadata(queryParam QueryParameter) []AppMetadataKey{
	return cobraSearch.getAppMetadataKeysByQuery(queryParam)
}

var parameterNames = []string{ "title", "version", "maintainerName", "maintainerEmail","company", "website", "source", "license"}

func (cobraSearch *CobraSearch) Init(){
	if cobraSearch.invertedIndex==nil{
		cobraSearch.invertedIndex = make(map[string]map[string]AppMetadataKeySet)
		for _, parameterName:= range parameterNames{
			cobraSearch.invertedIndex[parameterName] = make(map[string]AppMetadataKeySet)
		}
		log.Println("CobraSearch Inverted Index is initialized")
	}else{
		log.Println("CobraSearch Inverted index is already initialized.")
	}
}

func indexProperty(cobraSearch *CobraSearch, appMetadata AppMetadata, propertyName string , propertyValue string ){
	log.Printf("Indexing %s = %s for app metadata title: %s, version: %s\n", propertyName, propertyValue, appMetadata.Title, appMetadata.Version)
	var propertyIndex = cobraSearch.invertedIndex[propertyName]
	propertyValue = strings.TrimSpace(propertyValue)
	var key = AppMetadataKey{appMetadata.Title, appMetadata.Version}
	current := 	propertyIndex[propertyValue]
	propertyIndex[propertyValue] = current.Add(key)
}

func (cobraSearch *CobraSearch) getAppMetadataKeysByQuery(queryParam QueryParameter) []AppMetadataKey{
	keysQueriedFromParameters := [][]AppMetadataKey{}
	for _,parameterName := range parameterNames{
		queriedKeys := cobraSearch.getAppMetadataByProperty(queryParam, parameterName)
		if queriedKeys !=nil{
			keysQueriedFromParameters = append(keysQueriedFromParameters, queriedKeys)
		}
	}
	return IntersectAll(keysQueriedFromParameters)
}

func (cobraSearch *CobraSearch) getAppMetadataByProperty(queryParam QueryParameter, propertyName string) []AppMetadataKey{
	if queriedValue :=  getQueriedValue(queryParam, propertyName); queriedValue !=""{
		keysQueried := cobraSearch.invertedIndex[propertyName][queriedValue]
		return keysQueried.GetAllAppMetadataKeys()
	}else {
		return nil
	}
}

func getQueriedValue(parameter QueryParameter, queriedProperty string) string{
	switch queriedProperty {
	case "title":
		return strings.TrimSpace(parameter.Title)
	case "version":
		return strings.TrimSpace(parameter.Version)
	case "maintainerName":
		return strings.TrimSpace(parameter.MaintainerName)
	case "maintainerEmail":
		return strings.TrimSpace(parameter.MaintainerEmail)
	case "company":
		return strings.TrimSpace(parameter.Company)
	case "website":
		return strings.TrimSpace(parameter.Website)
	case "source":
		return strings.TrimSpace(parameter.Source)
	case "license":
		return strings.TrimSpace(parameter.License)
	default:
		return ""
	}
}