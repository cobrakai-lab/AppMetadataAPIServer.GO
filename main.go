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

const MetadataEndpoint = "v1/metadata"
const Port = "8888"

var cobraDB storage.Database = new(storage.CobraDB)

func main() {
	log.Println("Starting server")
	initServer().Run("0.0.0.0:" + Port)
}

func initServer() *gin.Engine{
	cobraDB.Init()
	router := gin.Default()
	router.GET(MetadataEndpoint, queryMetadata)
	router.POST(MetadataEndpoint, postMetadata)
	router.GET(MetadataEndpoint+"/:title/:version", getMetadataByTitleVersion)
	return router
}

func getMetadataByTitleVersion(c *gin.Context) {
	defer recovery(c)
	title := c.Param("title")
	version := c.Param("version")
	var key = storage.AppMetadataKey{title, version}
	var result = cobraDB.GetBulk([]storage.AppMetadataKey{key})
	if len(result) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No metadata found for title: %s, version: %s", title, version)})
	} else {
		c.IndentedJSON(http.StatusOK, result[0])
	}
}

func queryMetadata(c *gin.Context) {
	defer recovery(c)
	parameters:=storage.QueryParameter{
		Title:           c.Query("title"),
		Version:         c.Query("version"),
		MaintainerName:  c.Query("maintainerName"),
		MaintainerEmail: c.Query("maintainerEmail"),
		Company:         c.Query("company"),
		Website:         c.Query("website"),
		Source:          c.Query("source"),
		License:         c.Query("license"),
	}
	result:=cobraDB.Query(parameters)
	c.IndentedJSON(http.StatusOK, result)


}
func postMetadata(c *gin.Context) {
	defer recovery(c)
	var newMetadata AppMetadata
	if err := c.BindYAML(&newMetadata); err != nil {
		log.Printf("Something wrong with YAML format: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad format"})
		return
	}
	if err := validator.Validate(newMetadata); err != nil {
		log.Printf("Provided metadata is not valid. %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s", err)})
		return
	}
	if err := cobraDB.Create(newMetadata); err != nil {
		log.Printf("Error when writing to database: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s", err)})
		return
	}
	log.Printf("New metadata created successfully")
	c.IndentedJSON(http.StatusCreated, newMetadata)
}

func recovery(c *gin.Context){
	if err:=recover();err!=nil{
		log.Printf("Runtime error: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Something wrong, please try again later."})
	}
}