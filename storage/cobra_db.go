package storage

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"fmt"
)

func CreateKey(appMetadata AppMetadata) AppMetadataKey {
	return AppMetadataKey{
		Title:   appMetadata.Title,
		Version: appMetadata.Version,
	}
}

func Create(metadata AppMetadata) error{
	key:= CreateKey(metadata)
	if _, found :=dataCore[key]; found {
		return errors.New("duplicate record exists")
	}else{
		dataCore[key] = metadata
	}
	return nil
}

func GetBulk(keys []AppMetadataKey) []AppMetadata {
	var result []AppMetadata = []AppMetadata{}
	for _, key := range keys{
		if data, found :=dataCore[key]; found {
			fmt.Printf("data: %s", data.Title)
			result = append(result, data)
		}
	}
	fmt.Printf("Result: %s", result)
	return result
}

var dataCore = make(map[AppMetadataKey]AppMetadata)

type AppMetadataKey struct {
	Title   string
	Version string
}
