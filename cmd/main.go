package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhexists/dns/server"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	
	done := make(chan bool)
	
	go server.StartDNSServer(done)
	
	<-done
	
	server.StartHTTPServer()
}
