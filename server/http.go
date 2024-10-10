package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/controllers"
	. "github.com/shubhexists/dns/internal/logger"
)

// This is the start of HTTP Server that will expose out DB CRUD Operations
func StartHTTPServer() {
	router := gin.Default()

	records := router.Group("/records")
	{
		records.GET("/getRecords", controllers.GetRecordsByDomainID)
		records.POST("/createDomain", controllers.CreateDomain)
		records.POST("/createRecords", controllers.CreateRecord)
		records.DELETE("/deleteDomainByID", controllers.DeleteDomainByID)
		records.DELETE("/deleteDomain", controllers.DeleteDomainByID)
		records.PUT("/updateRecord", controllers.UpdateRecord)
	}

	// PORT Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	Log.Printf("HTTP server is running on port %s...\n", port)
	err := router.Run(":" + port)
	if err != nil {
		Log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
