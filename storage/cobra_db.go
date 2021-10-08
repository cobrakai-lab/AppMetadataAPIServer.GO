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
	Init()
}

type CobraDB struct {
	dataCore      map[AppMetadataKey]AppMetadata
}

func (cobraDb *CobraDB) Init(){
	cobraDb.dataCore = make(map[AppMetadataKey]AppMetadata)
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
		//todo add to inverted index
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

type AppMetadataKey struct {
	Title   string
	Version string
}
