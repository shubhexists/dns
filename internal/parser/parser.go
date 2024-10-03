package parser

import (
	"encoding/binary"

	"github.com/shubhexists/dns/internal/helpers"
	"github.com/shubhexists/dns/models"
)

func ParseDNSHeader(data []byte) models.DNSHeader {
	// `data` is 12 bytes.
	// First 2 bytes i.e. 1 and 2 (16 bits) is packedID
	packetID := binary.BigEndian.Uint16(data[0:2])

	// Byte No. 3 contains - QR, OPCode, AA and TC
	qrOpcodeAaTcRd := data[2]
	// qrOpcodeAaTcRd >> 7 shifts 7 places to right
	// So, 1101011 -> 7 right shifts -> 00000001
	// Take bitwise AND of 00000001 and 0x1 (00000001)
	// We will get the last bit
	qr := (qrOpcodeAaTcRd >> 7) & 0x1
	// Right shift 3 places and AND by 00001111 to get 4 bits
	opcode := (qrOpcodeAaTcRd >> 3) & 0xF
	// Right shift 2 places and AND by 00000001 to get next 1 bit
	aa := (qrOpcodeAaTcRd >> 2) & 0x1
	// Right shift by 1 place and AND
	tc := (qrOpcodeAaTcRd >> 1) & 0x1
	// get last bit
	rd := qrOpcodeAaTcRd & 0x1

	// Byte No. 4 contains - RD, AA and Z
	raZRcode := data[3]
	// 1 bit
	ra := (raZRcode >> 7) & 0x1
	// 0x7 is 00000111 -> get the last 3 bits
	z := (raZRcode >> 4) & 0x7
	// Remaining 4 bits
	rcode := raZRcode & 0xF

	// Byte 5 and 6 (16 bits) is Number of Questions
	qdcount := binary.BigEndian.Uint16(data[4:6])
	// Byte 7 and 8 (16 bits) is the Number of Answers
	ancount := binary.BigEndian.Uint16(data[6:8])
	// Byte 9 and 10 (16 bits) is Number of Authority
	nscount := binary.BigEndian.Uint16(data[8:10])
	// Byte 11 and 12 (16 bits) is the number of Additional Records
	arcount := binary.BigEndian.Uint16(data[10:12])

	return models.DNSHeader{
		PacketID: packetID,
		QR:       qr,
		Opcode:   opcode,
		AA:       aa,
		TC:       tc,
		RD:       rd,
		RA:       ra,
		Z:        z,
		RCode:    rcode,
		QDCount:  qdcount,
		ANCount:  ancount,
		NSCount:  nscount,
		ARCount:  arcount,
	}
}

func ParseDNSQuestion(data []byte) (models.DNSQuestion, int, int) {
	i := 12
	qnamebyte := [][]byte{}
	for data[i] != 0x00 {
		skip := int(data[i])
		i++
		qnamebyte = append(qnamebyte, data[i:i+skip])
		i = i + skip
	}
	i++

	var qname []string
	for _, b := range qnamebyte {
		qname = append(qname, helpers.ByteToString(b))
	}

	qtype := binary.BigEndian.Uint16(data[i : i+2])
	qclass := binary.BigEndian.Uint16(data[i+2 : i+4])

	return models.DNSQuestion{
		QName:  qname,
		QType:  qtype,
		QClass: qclass,
	}, i + 4, i - 12
}

func ParseDNSAnswer(data []byte) models.DNSAnswer {
	question, i, len := ParseDNSQuestion(data)

	name := question.QName
	index := i + len
	record_type := binary.BigEndian.Uint16(data[index : index+2])
	class := binary.BigEndian.Uint16(data[index+2 : index+4])
	ttl := binary.BigEndian.Uint32(data[index+4 : index+8])
	rdlength := binary.BigEndian.Uint16(data[index+8 : index+10])
	rddata := data[index+10 : index+10+int(rdlength)]

	return models.DNSAnswer{
		Name:     name,
		Type:     record_type,
		Class:    class,
		TTL:      ttl,
		RDLENGTH: rdlength,
		RDATA:    helpers.ByteToInt(rddata),
	}
}
