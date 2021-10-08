package validator

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"strings"
)

func Validate(metadata AppMetadata) error{
	if strings.TrimSpace(metadata.Title) ==""{
		return errors.New("Title is required.")
	}

	if strings.TrimSpace(metadata.Version) == ""{
		return errors.New("Version is required.")
	}
	//todo check other fields.

	return nil
}