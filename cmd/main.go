package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
	"github.com/shubhexists/dns/server"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = database.ConnectToDB()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = database.DB.AutoMigrate(&models.DNSRecords{})
	if err != nil {
		log.Fatal("Error migrating to new schema")
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	done := make(chan bool)

	go server.StartDNSServer(done)

	<-done

	server.StartHTTPServer()
}
