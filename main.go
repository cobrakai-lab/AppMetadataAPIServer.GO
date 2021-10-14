package main

import (
	. "AppMetadataAPIServerGo/model"
	"AppMetadataAPIServerGo/storage"
	"AppMetadataAPIServerGo/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const MetadataEndpoint = "v1/metadata"
const Port = "8888"

var cobraDB = storage.CobraDB{}
var serverStartTime = time.Now()
var getMetadataAPICount = 0
var postMetadataAPICount=0
var queryMetadataAPICount = 0
var deleteMetadataAPICount = 0
var statsViewed = 0

func main() {
	log.Println("Starting server")
	router:=initServer()
	preloadMockData()
	router.Run("0.0.0.0:" + Port)
}

func initServer() *gin.Engine{
	var cobraSearch storage.SearchEngine = new(storage.CobraSearch)
	cobraSearch.Init()
	cobraDB.Init(cobraSearch)
	router := gin.Default()
	router.GET(MetadataEndpoint, queryMetadata)
	router.POST(MetadataEndpoint, postMetadata)
	router.GET(MetadataEndpoint+"/:title/:version", getMetadataByTitleVersion)
	router.DELETE(MetadataEndpoint+"/:title/:version", deleteMetadata)
	router.GET(MetadataEndpoint+"/_stats", getServerStats)
	return router
}

func getMetadataByTitleVersion(c *gin.Context) {
	defer recovery(c)
	getMetadataAPICount+=1

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
	queryMetadataAPICount+=1
	const DefaultPageSize = 20
	pageNumber,err := strconv.Atoi(c.Query("page"))
	if err!=nil || pageNumber<1{
		pageNumber=1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err!=nil  || pageSize<0{
		pageSize = DefaultPageSize
	}
	parameters:=storage.QueryParameter{
		Title:           c.Query("title"),
		Version:         c.Query("version"),
		MaintainerName:  c.Query("maintainerName"),
		MaintainerEmail: c.Query("maintainerEmail"),
		Company:         c.Query("company"),
		Website:         c.Query("website"),
		Source:          c.Query("source"),
		License:         c.Query("license"),
		Page: 			 pageNumber,
		PageSize: 		 pageSize,
	}

	if hasUnrecognizedQuery(c.Request.URL.Query()){
		c.IndentedJSON(http.StatusOK, []AppMetadata{})
		return
	}
	result:=cobraDB.Query(parameters)
	c.IndentedJSON(http.StatusOK, result)
}

func hasUnrecognizedQuery(query url.Values) bool {
	recognizedQueryCount := 0
	for _, parameterName := range storage.ParameterNames{
		if _,found:= query[parameterName];found{
			recognizedQueryCount++
		}
	}
	if _, found:=query["page"];found{
		recognizedQueryCount++
	}
	if _, found:=query["pageSize"];found{
		recognizedQueryCount++
	}
	log.Printf("recognized query count: %d, actual: %d", recognizedQueryCount, len(query))
	return len(query) != recognizedQueryCount
}


func postMetadata(c *gin.Context) {
	defer recovery(c)
	postMetadataAPICount+=1

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

func deleteMetadata(c *gin.Context){
	defer recovery(c)
	deleteMetadataAPICount+=1

	title := c.Param("title")
	version := c.Param("version")
	var key = storage.AppMetadataKey{title, version}
	deleted, err := cobraDB.Delete(key)
	if err==nil{
		c.IndentedJSON(http.StatusOK, deleted)
	}else{
		c.IndentedJSON(http.StatusOK, gin.H{"result":"metadata does not exist"})
	}
}

func getServerStats(c *gin.Context){
	defer recovery(c)
	statsViewed+=1

	c.IndentedJSON(http.StatusOK, gin.H{
		"server_start_time":fmt.Sprintf(serverStartTime.String()),
		"get_metadata_api_called_counter": getMetadataAPICount,
		"query_metadata_api_called_counter": queryMetadataAPICount,
		"post_metadata_api_called_counter": postMetadataAPICount,
		"delete_metadata_api_called_counter": deleteMetadataAPICount,
		"stats_viewed": statsViewed,
	})
}

func recovery(c *gin.Context){
	if err:=recover();err!=nil{
		log.Printf("Runtime error: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"Something wrong, please try again later."})
	}
}


func preloadMockData(){
	mockData := []AppMetadata{
		{
			Title: "mock app",
			Version: "1.0",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
			},
			Company:     "Cobrakai",
			Website:     "www.somewhere.com",
			Source:      "http://github.com",
			License:     "MIT",
			Description: "Looks legit!",
		},
		{
			Title: "mock app",
			Version: "1.1",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
			},
			Company:     "Cobrakai",
			Website:     "www.somewhere.com",
			Source:      "http://github.com",
			License:     "MIT",
			Description: "Looks legit!",
		},
		{
			Title: "real app",
			Version: "1.0",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
			},
			Company:     "Cobrakai",
			Website:     "www.somewhere.com",
			Source:      "http://github.com",
			License:     "Apache-2.0",
			Description: "Looks legit!",
		},
		{
			Title: "real app2",
			Version: "1.0",
			Maintainers: []Maintainer{
				{"kai", "i@g.com"},
				{"cobra", "a@b.com"},
			},
			Company:     "Cobrakai",
			Website:     "www.somewhere.com",
			Source:      "http://github.com",
			License:     "BSD",
			Description: "Looks legit!",
		},
	}

	for _,metadata:= range(mockData){
		cobraDB.Create(metadata)
	}
}
