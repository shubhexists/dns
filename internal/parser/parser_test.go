package parser

import (
	"bytes"
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

func TestParseDNSQuestion(t *testing.T) {
	data := []byte{
		0xdb, 0x42,
		0x01,
		0x00,
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x03, 0x77,
		0x77, 0x77,
		0x06, 0x67,
		0x6f, 0x6f,
		0x67, 0x6c,
		0x65, 0x03,
		0x63, 0x6f,
		0x6d, 0x00,
		0x00, 0x01,
		0x00, 0x01,
	}

	expected := models.DNSQuestion{
		QName:  []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00},
		QType:  1,
		QClass: 1,
	}

	result := ParseDNSQuestion(data)

	fmt.Printf("%+v\n", result)
	// Compare the QName byte slice using bytes.Equal
	if !bytes.Equal(result.QName, expected.QName) {
		t.Errorf("ParseDNSQuestion() QName = %+v, want %+v", result.QName, expected.QName)
	}

	// Compare QType and QClass directly
	if result.QType != expected.QType {
		t.Errorf("ParseDNSQuestion() QType = %+v, want %+v", result.QType, expected.QType)
	}

	if result.QClass != expected.QClass {
		t.Errorf("ParseDNSQuestion() QClass = %+v, want %+v", result.QClass, expected.QClass)
	}
}
