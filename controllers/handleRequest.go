package controllers

import (
	"net"

	"github.com/shubhexists/dns/internal/helpers"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/internal/parser"
)

func HandleDNSRequest(conn *net.UDPConn) {
	buffer := make([]byte, 512)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		Log.Printf("Failed to read from UDP connection: %v", err)
		return
	}
	dnsHeader := parser.ParseDNSHeader(buffer[:12])
	dnsQuestion, _, _ := parser.ParseDNSQuestion(buffer[12:n])

	Log.Printf("Received DNS query from %s: %+v", addr, dnsQuestion)

	response := helpers.BuildDNSResponse(dnsHeader, dnsQuestion)

	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		Log.Printf("Failed to write to UDP connection: %v", err)
		return
	}

	Log.Printf("Sent response to %s", addr)
}
