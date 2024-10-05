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

	router.GET("/getRecords", controllers.GetRecord)
	router.POST("/createRecords", controllers.CreateRecord)
	router.PUT("/updateRecord", controllers.UpdateRecord)
	router.DELETE("/deleteRecord", controllers.DeleteRecord)
	
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
