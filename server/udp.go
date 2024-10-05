package server

import (
	"fmt"
	"log"
	"net"

	"github.com/shubhexists/dns/controllers"
)

// DNS Server start
func StartDNSServer(done chan bool) {
	addr := net.UDPAddr{
		Port: 5350,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP port 5350: %v", err)
	}
	defer conn.Close()

	fmt.Println("DNS server is running on port 5350...")
	done <- true

	for {
		controllers.HandleDNSRequest(conn)
	}
}
