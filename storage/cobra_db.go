package storage

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"log"
	"sort"
	"strings"
	"sync"
)

type Database interface{
	Create(metadata AppMetadata) error
	GetBulk(keys []AppMetadataKey) []AppMetadata
	Query(parameter QueryParameter) []AppMetadata
	Delete(key AppMetadataKey) (AppMetadata, error)
	Init()
}

type CobraDB struct {
	dataCore    map[AppMetadataKey]AppMetadata
	cobraSearch SearchEngine
	mutex sync.Mutex
}

func (cobraDb *CobraDB) Init(cobraSearch SearchEngine){
	cobraDb.dataCore = make(map[AppMetadataKey]AppMetadata)
	cobraDb.cobraSearch = cobraSearch
	log.Println("CobraDB is initialized.")
}

func (cobraDb *CobraDB) Create(metadata AppMetadata) error{
	cobraDb.mutex.Lock()
	defer cobraDb.mutex.Unlock()
	key:= AppMetadataKey{ metadata.Title, metadata.Version}
	if _, found := cobraDb.dataCore[key]; found {
		return errors.New("duplicate record exists")
	}else{
		cobraDb.dataCore[key] = metadata
		cobraDb.cobraSearch.IndexMetadata(metadata)
	}
	return nil
}

func (cobraDb *CobraDB) GetBulk(keys []AppMetadataKey) []AppMetadata {
	var result = []AppMetadata{}
	for _, key := range keys{
		if data, found := cobraDb.dataCore[key]; found {
			result = append(result, data)
		}
	}
	return result
}

func (cobraDb *CobraDB) Query(parameter QueryParameter) []AppMetadata{
	log.Printf("Querying with parameter: %+v", parameter)
	keys := []AppMetadataKey{}
	if isEmptyQuery(parameter){
		keys = cobraDb.getAllKeys()
	}else{
		keys = cobraDb.cobraSearch.QueryMetadata(parameter)
	}
	sorted(keys)
	keys = getKeysByPage(parameter.Page, parameter.PageSize,keys)
	return cobraDb.GetBulk(keys)
}

func (cobraDb *CobraDB) Delete(key AppMetadataKey) (AppMetadata, error){
	if metadata, found := cobraDb.dataCore[key];found{
		delete(cobraDb.dataCore, key)
		cobraDb.cobraSearch.Delete(key)
		return metadata, nil
	}
	return AppMetadata{}, errors.New("No metadata found")
}

func (cobraDb *CobraDB) getAllKeys() []AppMetadataKey{
	keys := []AppMetadataKey{}
	for key, _ := range cobraDb.dataCore{
		keys = append(keys, key)
	}
	return keys
}

func isEmptyQuery(parameter QueryParameter) bool {
	return strings.TrimSpace(
		parameter.Title+
		parameter.Version+
		parameter.MaintainerEmail+
		parameter.MaintainerName+
		parameter.Website+
		parameter.Company+
		parameter.Source+
		parameter.License) == ""
}

func getKeysByPage(pageNumber int, pageSize int, keys []AppMetadataKey) []AppMetadataKey{
	startIndex := (pageNumber-1)*pageSize
	if startIndex>=len(keys){
		return []AppMetadataKey{}
	}
	exclusiveEndIndex := startIndex+pageSize
	if exclusiveEndIndex>len(keys){
		exclusiveEndIndex = len(keys)
	}
	return keys[startIndex: exclusiveEndIndex]
}

func sorted(keys []AppMetadataKey){
	sort.SliceStable(keys, func(i, j int) bool {
		if keys[i].Title < keys[j].Title{
			return true
		}
		if keys[i].Title>keys[j].Title{
			return false
		}
		return keys[i].Version< keys[j].Version
	})
}


type AppMetadataKey struct {
	Title   string
	Version string
}
