package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// This is the start of HTTP Server that will expose out DB CRUD Operations
func StartHTTPServer() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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
