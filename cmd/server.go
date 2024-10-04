package main

import (
	"fmt"
	"log"
	"net"

	"github.com/shubhexists/dns/controller"
)

func main() {
	addr := net.UDPAddr{
		// We can not use 53 as most of the Linux machines have 53 blocked
		// Probably we'll see if we can take this from a config or env file
		Port: 5350,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP port 5350: %v", err)
	}

	defer conn.Close()

	fmt.Println("DNS server is running on port 5350...")

	for {
		controller.HandleRequest(conn)
	}
}
