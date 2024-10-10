package server

import (
	"net"

	"github.com/shubhexists/dns/controllers"
	. "github.com/shubhexists/dns/internal/logger"
)

// DNS Server start
func StartDNSServer(done chan bool) {
	addr := net.UDPAddr{
		Port: 5350,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		Log.Fatalf("Failed to listen on UDP port 5350: %v", err)
	}
	defer conn.Close()

	Log.Println("DNS server is running on port 5350...")
	done <- true

	for {
		controllers.HandleDNSRequest(conn)
	}
}
