package main

import (
	. "AppMetadataAPIServerGo/model"
	"AppMetadataAPIServerGo/storage"
	"AppMetadataAPIServerGo/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)



var AppMetadatas []AppMetadata = []AppMetadata{
	{Title: "app1", Version: "1.0", Maintainers: []Maintainer{
		{
		Name:  "kai",
		Email: "i@g.com",
		}}, Company: "Cobrakai", Website: "www.abc.com", Source: "github", License: "mit", Description: "boom!"},
}

const MetadataEndpoint = "v1/metadata"
const Port = "8888"

func main() {
	log.Println("Starting server")
	router := gin.Default()
	router.GET(MetadataEndpoint, getMetadata)
	router.POST(MetadataEndpoint, postMetadata)
	router.Run("localhost:"+Port)
}
func getMetadata(c *gin.Context) {

	title := c.Query("title")
	version := c.Query("version")
	log.Printf("title: %s, version: %s", title, version)
	key := storage.AppMetadataKey{
		title,
	version,
	}
	c.IndentedJSON(http.StatusOK, storage.GetBulk([]storage.AppMetadataKey{key}))
}
func postMetadata(c *gin.Context) {
	var newMetadata AppMetadata

	if err := c.BindYAML(&newMetadata); err != nil {
		log.Printf("Something wrong with binding YAML: %s", err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Bad format"})
		return
	}
	if err:= validator.Validate(newMetadata); err!=nil {
		log.Printf("Provide metadata is not valid. %s", err)
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":fmt.Sprintf("%s", err)})
		return
	}
	if err:= storage.Create(newMetadata); err!=nil{
		log.Printf("Error when writing to database: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s", err)})
		return
	}
	log.Printf("New metadata created successfully")
	c.IndentedJSON(http.StatusCreated, newMetadata)
}

