package storage

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"log"
)

type Database interface{
	CreateKey(metadata AppMetadata) AppMetadataKey
	Create(metadata AppMetadata) error
	GetBulk(keys []AppMetadataKey) []AppMetadata
	Query(parameter QueryParameter) []AppMetadata
	Init()
}

type CobraDB struct {
	dataCore      map[AppMetadataKey]AppMetadata
	cobraSearch   cobraSearch
}

func (cobraDb *CobraDB) Init(){
	cobraDb.dataCore = make(map[AppMetadataKey]AppMetadata)
	cobraDb.cobraSearch = cobraSearch{}
	cobraDb.cobraSearch.initInvertedIndex()
	log.Println("CobraDB is initialized.")
}

func (cobraDb *CobraDB) CreateKey(appMetadata AppMetadata) AppMetadataKey {
	return AppMetadataKey{
		Title:   appMetadata.Title,
		Version: appMetadata.Version,
	}
}

func (cobraDb *CobraDB) Create(metadata AppMetadata) error{
	key:= cobraDb.CreateKey(metadata)
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
	var keys = cobraDb.cobraSearch.QueryMetadata(parameter)
	return cobraDb.GetBulk(keys)
}

type AppMetadataKey struct {
	Title   string
	Version string
}
