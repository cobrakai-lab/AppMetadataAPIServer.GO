package storage

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"log"
	"sync"
)

type Database interface{
	Create(metadata AppMetadata) error
	GetBulk(keys []AppMetadataKey) []AppMetadata
	Query(parameter QueryParameter) []AppMetadata
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
	var keys = cobraDb.cobraSearch.QueryMetadata(parameter)
	return cobraDb.GetBulk(keys)
}

type AppMetadataKey struct {
	Title   string
	Version string
}
