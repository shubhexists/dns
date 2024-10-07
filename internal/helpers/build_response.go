package helpers

import (
	"encoding/binary"
	"strings"

	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
)

func BuildDNSResponse(header models.DNSHeader, question models.DNSQuestion) []byte {
	response := make([]byte, 512)

	offset := 0
	binary.BigEndian.PutUint16(response[offset:], header.PacketID)
	offset += 2

	//-------------------------------------------------------------------------------------------------

	// QR | Opcode | AA | TC | RD
	// NOTE - AA should always be ) for the query. We'll check for that domain on our DB and check if we have it and set AA to 1
	// TO SEE - How to handle truncation
	flags := uint16(1)<<15 | uint16(header.Opcode)<<11 | uint16(header.AA)<<10 |
		uint16(header.TC)<<9 | uint16(header.RD)<<8

	// Bitwise OR
	// RA | Z | RCode
	// To see - DO We support Recursion lol?
	flags |= uint16(header.RA)<<7 | uint16(header.Z)<<4 | uint16(header.RCode)
	binary.BigEndian.PutUint16(response[offset:], flags)
	offset += 2

	binary.BigEndian.PutUint16(response[offset:], header.QDCount)
	offset += 2
	// number of answers ( Can we send more than 1 answers? )
	binary.BigEndian.PutUint16(response[offset:], header.ANCount)
	offset += 2
	// probably there might be NS Records
	binary.BigEndian.PutUint16(response[offset:], header.NSCount)
	offset += 2
	// NOT SURE WHAT TO SET HERE
	binary.BigEndian.PutUint16(response[offset:], header.ARCount)
	offset += 2

	//--------------------------------------------------------------------------------------------------

	// THIS PART SETS QNAME, IS GOOD ACC TO ME
	for _, label := range question.QName {
		response[offset] = byte(len(label))
		offset++
		copy(response[offset:], label)
		offset += len(label)
	}
	response[offset] = 0
	offset++

	// GOOD
	binary.BigEndian.PutUint16(response[offset:], uint16(question.QType))
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], question.QClass)
	offset += 2

	//--------------------------------------------------------------------------------------------------//
	//                                          ANSWER SECTION                                          //
	//--------------------------------------------------------------------------------------------------//

	var name, baseURL string
	if len(question.QName) > 2 {
		name = question.QName[0]
		baseURL = strings.Join(question.QName[1:], ".")
	} else {
		name = ""
		baseURL = strings.Join(question.QName, ".")
	}

	var record []models.DNSRecords
	if err := database.DB.Where("name = ? AND base_url = ?", name, baseURL).Find(&record).Error; err != nil {
		// Figure out how to send errors :D
	}

	switch question.QType {
	case models.QTYPE_A:

	case models.QTYPE_AAAA:

	case models.QTYPE_CNAME:

	case models.QTYPE_MX:

	case models.QTYPE_NS:

	case models.QTYPE_PTR:

	case models.QTYPE_SOA:

	case models.QTYPE_TXT:

	default:

	}

	// Authority section

	// Additional Section

	return response
}
