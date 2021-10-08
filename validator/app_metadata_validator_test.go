package validator_test

import (
	"AppMetadataAPIServerGo/model"
	"AppMetadataAPIServerGo/validator"
	"testing"
)

func TestAllFieldsRequired(t *testing.T) {
	var appMetadata = createValidAppMetadata()
	if validator.Validate(appMetadata)!=nil {
		t.Error("Validator failed to validate a valid app metadata.")
	}

	appMetadata.Title = ""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when title is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Version=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when version is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Maintainers=[]model.Maintainer{}
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when maintainer is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Maintainers[0].Name=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when maintainer name is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Maintainers[0].Email=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when maintainer email is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Company=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when company is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Website=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when website is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Source=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when source is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.License=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when license is empty")
	}

	appMetadata = createValidAppMetadata()
	appMetadata.Description=""
	if validator.Validate(appMetadata)==nil {
		t.Error("Validator failed to validate when description is empty")
	}
}

func TestValidateEmails(t *testing.T) {
	validEmails := []string{
		"a@b.com",
		"i@g.co",
		"a+b_123-ip@xyz.jkl.app",
		"firstname.lastname@example.com",
		"firstname+lastname@example.com",
		"1234567890@example.com",
		"email@example-one.com",
		"_______@example.com",
	}

	invalidEmails := []string{
		"#@%^%#$@#$@#.com",
		"@example.com",
		"email.example.com",
		"email@example@example.com",
		".email@example.com",
		"email.@example.com",
		"email..email@example.com",
		"email@example..com",
		"Abc..123@example.com",
	}

	for _, invalidEmail:= range invalidEmails{
		var appMetadata = createValidAppMetadata()
		appMetadata.Maintainers[0].Email=invalidEmail
		if validator.Validate(appMetadata)==nil {
			t.Errorf("Validator failed to validate invalid email: %s", invalidEmail)
		}
	}

	for _, validEmail := range validEmails{
		var appMetadata = createValidAppMetadata()
		appMetadata.Maintainers[0].Email=validEmail
		if validator.Validate(appMetadata)!=nil {
			t.Errorf("Validator failed to validate correct email: %s", validEmail)
		}
	}
}

func createValidAppMetadata() model.AppMetadata{
	return model.AppMetadata{
		Title:   "app1",
		Version: "1.0",
		Maintainers: []model.Maintainer{
			{
				Name:  "kai",
				Email: "i@g.com",
			},
		},
		Company: "Cobrakai",
		Website: "www.abc.com",
		Source: "github",
		License: "mit",
		Description: "boom!"}
}
