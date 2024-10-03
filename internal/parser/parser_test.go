package parser

import (
	"fmt"
	"reflect"
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
		0x12, 0x34, // Packet ID = 0x1234
		0x85,       // QR = 1, Opcode = 0, AA = 1, TC = 0, RD = 1
		0x80,       // RA = 1, Z = 0, RCode = 0
		0x00, 0x01, // QDCount = 1
		0x00, 0x02, // ANCount = 2
		0x00, 0x03, // NSCount = 3
		0x00, 0x04, // ARCount = 4
		0x03, 0x77, // 3 length
		0x77, 0x77,
		0x06, 0x67, // 6 length
		0x6f, 0x6f,
		0x67, 0x6c,
		0x65, 0x03, //3 length
		0x63, 0x6f,
		0x6d, 0x00, // so its www.google.com and then null character
		0x00, 0x01, // QType = 1
		0x00, 0x01, //QClass = 1
	}

	expected := models.DNSQuestion{
		QName:  []string{"www", "google", "com"},
		QType:  1,
		QClass: 1,
	}

	result, _, _ := ParseDNSQuestion(data)

	fmt.Printf("%+v\n", result)
	// Compare the QName byte slice using bytes.Equal
	if !reflect.DeepEqual(result.QName, expected.QName) {
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

func TestParseDNSAnswer(t *testing.T) {
	data := []byte{
		0x12, 0x34, // Packet ID = 0x1234
		0x85,       // QR = 1, Opcode = 0, AA = 1, TC = 0, RD = 1
		0x80,       // RA = 1, Z = 0, RCode = 0
		0x00, 0x01, // QDCount = 1
		0x00, 0x02, // ANCount = 2
		0x00, 0x03, // NSCount = 3
		0x00, 0x04, // ARCount = 4
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
		0x00, 0x00, 0x01, 0x2c,
		0x00, 0x04,
		0x8E, 0xFA, 0x48, 0x64,
	}

	expected := models.DNSAnswer{
		Name:     []string{"www", "google", "com"},
		Type:     1,
		Class:    1,
		TTL:      300,
		RDLENGTH: 4,
		RDATA:    []int{142, 250, 72, 100},
	}

	result := ParseDNSAnswer(data)
	fmt.Printf("%+v\n", result)

	if !reflect.DeepEqual(result.Name, expected.Name) {
		t.Errorf("Name mismatch. Got: %+v, Want: %+v", result.Name, expected.Name)
	}

	// Check other fields
	if result.Type != expected.Type {
		t.Errorf("Type mismatch. Got: %d, Want: %d", result.Type, expected.Type)
	}

	if result.Class != expected.Class {
		t.Errorf("Class mismatch. Got: %d, Want: %d", result.Class, expected.Class)
	}

	if result.TTL != expected.TTL {
		t.Errorf("TTL mismatch. Got: %d, Want: %d", result.TTL, expected.TTL)
	}

	if result.RDLENGTH != expected.RDLENGTH {
		t.Errorf("RDLENGTH mismatch. Got: %d, Want: %d", result.RDLENGTH, expected.RDLENGTH)
	}

	// Check RDATA field
	if !reflect.DeepEqual(result.RDATA, expected.RDATA) {
		t.Errorf("RDATA mismatch. Got: %+v, Want: %+v", result.RDATA, expected.RDATA)
	}
}
