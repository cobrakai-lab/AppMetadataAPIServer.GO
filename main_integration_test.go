package main

import (
	. "AppMetadataAPIServerGo/model"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)
var validInputsDirectory =  "./_test_data/_valid_inputs"
var invalidInputsDirectory = "./_test_data/_invalid_inputs"
var allAppMetadatas = getAppMetadataFromPayloads(getPayloadsFromFile(validInputsDirectory))

func TestIntegration(t *testing.T){
	validPayloads := getPayloadsFromFile(validInputsDirectory)
	invalidPayloads := getPayloadsFromFile(invalidInputsDirectory)

	testServer := httptest.NewServer(initServer())
	defer testServer.Close()

	endpoint := fmt.Sprintf("%s/%s", testServer.URL, MetadataEndpoint)

	testPostValidData(t, validPayloads, endpoint)
	testPostInvalidData(t, invalidPayloads, endpoint)
	testGetByTitleVersionAPI(t, endpoint)
	testQueryAPI(t,endpoint)
	testDeleteAPI(t, validPayloads, endpoint)
}

func testQueryAPI(t *testing.T, endpoint string) {
	queries:= []string{
		"license=Apache-2.0",
		"title=mit",
		"version=1.0.1",
		"maintainerName=firstmaintainer+app1&website=https://website2.com",
		"wrongKey=nothing",
		"title=Valid+App+1",
	}

	expectedResults := [][]AppMetadata{
		allAppMetadatas,
		{},
		{allAppMetadatas[1], allAppMetadatas[2]},
		{allAppMetadatas[2] },
		{},
		{allAppMetadatas[0], allAppMetadatas[1] },

	}

	for i,query :=range queries{
		validateQueryResults(t, endpoint+"?"+query, expectedResults[i])
	}
}

func validateQueryResults(t *testing.T, queryURI string, expectedResult []AppMetadata){
	res,err:=http.Get(queryURI)
	if err!=nil{
		t.Fatalf("Something wrong with Query URI %s, Error: %s",queryURI, err)
	}
	assert.Equal(t, res.StatusCode, 200)
	actual:= getAppMetadataSliceFromResponse(res)
	assert.ElementsMatch(t, expectedResult, actual)
}

func testGetByTitleVersionAPI(t *testing.T, endpoint string) {
	queries:= []string{
		"/Valid App 1/0.0.1/",
		"/Valid App 1/1.0.1/",
		"/Valid App 2/1.0.1/",
		"/Valid App 2/1.0.2/",
	}

	for i, query:= range queries{
		res,err:=http.Get(endpoint+query)
		if err!=nil{
			t.Fatalf("Something wrong with integration test: %s", err)
		}
		if res.StatusCode!=200{
			t.Errorf("Failed to get response from %s. Expected status code = 200, actual: %d", query, res.StatusCode)
		}
		returnedAppMetadata := getAppMetadataFromResponse(res)
		expectedAppMetadata := allAppMetadatas[i]

		if !reflect.DeepEqual(expectedAppMetadata, returnedAppMetadata){
			t.Errorf("expected metatdata not equal returned metadata. Expected: %+v, returned: %+v", expectedAppMetadata, returnedAppMetadata)
		}
	}

	res,err:=http.Get(endpoint+"/not/there/")
	if err!=nil{
		t.Fatalf("Something wrong with integration test: %s", err)
	}
	if res.StatusCode!=404{
		t.Errorf("Expected status code = 404, actual: %d", res.StatusCode)
	}
}

func testPostInvalidData(t *testing.T, invalidPayloads []string, endpoint string) {
	for _, invalidPayload := range invalidPayloads {
		res, err := http.Post(endpoint, "text/plain", strings.NewReader(invalidPayload))
		if err != nil {
			t.Fatalf("Something wrong with integration test: %s", err)
		}
		if res.StatusCode != 400 {
			t.Errorf("Failed to reject invalid payloads.Expected status code = 400, actual: %d, payload:\n %s", res.StatusCode, invalidPayload)
		}
	}
}

func testPostValidData(t *testing.T, validPayloads []string, endpoint string) {
	for _, validPayload := range validPayloads {
		res, err := http.Post(endpoint, "text/plain", strings.NewReader(validPayload))
		if err != nil {
			t.Fatalf("Something wrong with integration test: %s", err)
		}
		if res.StatusCode != 201 {
			responseContent, _ := io.ReadAll(res.Body)
			t.Errorf("Expected status code = 201, actual: %d, body: %s", res.StatusCode, responseContent)
		}
	}
	res, err := http.Post(endpoint, "text/plain", strings.NewReader(validPayloads[0]))
	if err != nil {
		t.Fatalf("Something wrong with integration test: %s", err)
	}
	if res.StatusCode != 400 {
		t.Errorf("Failed to reject duplicate payloads.Expected status code = 400, actual: %d, payload:\n %s", res.StatusCode, validPayloads[0])
	}
}

func testDeleteAPI(t *testing.T, validPayloads []string, endpoint string) {
	for _, validPayload := range validPayloads {
		http.Post(endpoint, "text/plain", strings.NewReader(validPayload))
	}
	toDelete:="/Valid App 1/0.0.1/"
	notToDelete := []string{
		"/Valid App 1/1.0.1/",
		"/Valid App 2/1.0.1/",
		"/Valid App 2/1.0.2/",
	}
	request, _:=http.NewRequest(http.MethodDelete, endpoint+toDelete,nil)

	//before delete
	validateQueryResults(t, endpoint+"?title=Valid+App+1", []AppMetadata{allAppMetadatas[0],allAppMetadatas[1]})
	res,err:=http.Get(endpoint+toDelete)
	assert.Equal(t, getAppMetadataFromResponse(res), allAppMetadatas[0])

	//do delete
	client:=&http.Client{}
	res,err = client.Do(request)
	if err!=nil{
		t.Errorf("Failed to handle Delete request. Error: %s", err)
	}
	//verify
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, getAppMetadataFromResponse(res),allAppMetadatas[0])
	res,err=http.Get(endpoint+toDelete)
	assert.Equal(t, res.StatusCode, 404)

	for _, remains := range notToDelete{
		res,err=http.Get(endpoint+remains)
		assert.Equal(t, res.StatusCode, 200)
	}
	validateQueryResults(t, endpoint+"?title=Valid+App+1", []AppMetadata{allAppMetadatas[1]})
}

func getAppMetadataFromResponse(res *http.Response) AppMetadata{
	responseContent, _ := io.ReadAll(res.Body)
	appMetadata := AppMetadata{}
	json.Unmarshal(responseContent, &appMetadata)
	return appMetadata
}

func getAppMetadataSliceFromResponse(res *http.Response) []AppMetadata{
	responseContent, _ := io.ReadAll(res.Body)
	metadatas:= []AppMetadata{}
	json.Unmarshal(responseContent, &metadatas)
	return metadatas
}

func getPayloadsFromFile(directory string) []string{
	var payloads []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("walking: %s", path)
		if !info.IsDir(){
			content, err := os.ReadFile(path)
			if err!=nil{
				return err
			}
			payloads = append(payloads, string(content))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return payloads
}

func getAppMetadataFromPayloads(payloads []string) []AppMetadata{
	var appMetadatas []AppMetadata
	for _, payload := range payloads{
		parsed := AppMetadata{}
		yaml.Unmarshal([]byte(payload), &parsed)
		appMetadatas = append(appMetadatas, parsed)
	}
	return appMetadatas
}