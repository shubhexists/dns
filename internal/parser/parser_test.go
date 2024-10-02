package parser

import (
	"fmt"
	"testing"

	"github.com/shubhexists/dns/models"
)

func TestParseDNSHeader(t *testing.T) {
	data := []byte{
		0x12, 0x34, // Packet ID = 0x1234
		0x85,       // QR = 1, Opcode = 0, AA = 1, TC = 0, RD = 1
		0x80,       // RA = 1, Z = 0, RCode = 0
		0x00, 0x01, // QDCount = 1
		0x00, 0x02, // ANCount = 2
		0x00, 0x03, // NSCount = 3
		0x00, 0x04, // ARCount = 4
	}

	expected := models.DNSHeader{
		PacketID: 0x1234,
		QR:       1,
		Opcode:   0,
		AA:       1,
		TC:       0,
		RD:       1,
		RA:       1,
		Z:        0,
		RCode:    0,
		QDCount:  1,
		ANCount:  2,
		NSCount:  3,
		ARCount:  4,
	}

	result := ParseDNSHeader(data)

	fmt.Printf("%+v\n", result)
	if result != expected {
		t.Errorf("ParseDNSHeader() = %+v, want %+v", result, expected)
	}
}
