package validator

import (
	. "AppMetadataAPIServerGo/model"
	"errors"
	"net/mail"
	"strings"
)

func Validate(metadata AppMetadata) error {
	if strings.TrimSpace(metadata.Title) == "" {
		return errors.New("title is required")
	}

	if strings.TrimSpace(metadata.Version) == "" {
		return errors.New("version is required")
	}

	if len(metadata.Maintainers)==0{
		return errors.New("maintainers cannot be empty")
	}

	for _, maintainer := range metadata.Maintainers{
		if strings.TrimSpace(maintainer.Name) == "" {
			return errors.New( "maintainer name is required")
		}
		if strings.TrimSpace(maintainer.Email) == "" {
			return errors.New("maintainer email is required")
		}

		if !isEmailValid(strings.TrimSpace(maintainer.Email)) {
			return errors.New("email is not valid")
		}
	}

	if strings.TrimSpace(metadata.Company) == ""{
		return errors.New("company is required")
	}

	if strings.TrimSpace(metadata.Website) == ""{
		return errors.New("website is required")
	}

	if strings.TrimSpace(metadata.Source) == ""{
		return errors.New("source is required")
	}

	if strings.TrimSpace(metadata.License) == ""{
		return errors.New("license is required")
	}

	if strings.TrimSpace(metadata.Description) == ""{
		return errors.New("description is required")
	}

	return nil
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err==nil
}