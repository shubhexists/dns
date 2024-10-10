package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
	"github.com/shubhexists/dns/server"
)

func init() {
	InitializeLogger()
	err := godotenv.Load()

	if err != nil {
		Log.Fatal("Error loading .env file")
	}

	_, err = database.ConnectToDB()
	if err != nil {
		Log.Fatal("Error connecting to database")
	}

	err = database.DB.AutoMigrate(&models.SOARecord{})
	if err != nil {
		Log.Fatal("Error migrating SOARecord to new schema")
	}

	err = database.DB.AutoMigrate(&models.Domain{})
	if err != nil {
		Log.Fatal("Error migrating Domain to new schema")
	}

	err = database.DB.AutoMigrate(&models.DNSRecord{})
	if err != nil {
		Log.Fatal("Error migrating DNSRecord to new schema")
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	done := make(chan bool)

	go server.StartDNSServer(done)

	<-done

	server.StartHTTPServer()
}
