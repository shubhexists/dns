package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhexists/dns/controllers"
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

	sno, err := controllers.CheckForSOA()
	Log.Println("SNO: ", sno)

	if err != nil {
		soa := models.SOARecord{
			PrimaryNS:  "ns1.shubh.sh",
			AdminEmail: strings.Replace("shubh622005@gmail.com", "@", ".", 1),
			Serial:     sno,
			Refresh:    86400,
			Retry:      7200,
			Expire:     3600,
			TTL:        3600,
		}

		if err := database.DB.Create(&soa).Error; err != nil {
			Log.Fatal("Error creating SOA Records")
		}

		Log.Println("SOA records created successfully")
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	done := make(chan bool)

	go server.StartDNSServer(done)

	<-done

	server.StartHTTPServer()
}
