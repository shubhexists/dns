package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/controllers"
)

// This is the start of HTTP Server that will expose out DB CRUD Operations
func StartHTTPServer() {
	router := gin.Default()

	records := router.Group("/records")
	{
		records.GET("/getRecordById", controllers.GetRecordByID)
		records.GET("/getRecordsByName", controllers.GetRecordsByName)
		records.POST("/createRecords", controllers.CreateRecord)
		records.DELETE("/deleteRecordByID", controllers.DeleteRecordByID)
		records.DELETE("/deleteRecordsByName", controllers.DeleteRecordsByName)
		records.PUT("/updateRecord", controllers.UpdateRecordByID)
	}
	
	// PORT Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("HTTP server is running on port %s...\n", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
