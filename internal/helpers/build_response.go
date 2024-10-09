package helpers

import (
	"encoding/binary"

	responsehandlers "github.com/shubhexists/dns/internal/helpers/response_handlers"
	"github.com/shubhexists/dns/models"
)

const (
	DNSHeaderSize      = 12
	PointerCompression = 0xC00C
	IPv4AddressLength  = 4
	IPv6AddressLength  = 16
)

func BuildDNSResponse(header models.DNSHeader, question models.DNSQuestion) []byte {
	response := make([]byte, 512)
	offset := 0

	offset = writeHeader(response, header)

	offset = writeQuestion(response, offset, question)

	var err error
	switch question.QType {
	case models.QTYPE_A:
		offset, err = handleARecord(response, offset, question)
	case models.QTYPE_AAAA:
		offset, err = handleAAAARecord(response, offset, question)
	case models.QTYPE_CNAME:
		offset, err = handleCNAMERecord(response, offset, question)
	case models.QTYPE_MX:
		offset, err = handleMXRecord(response, offset, question)
	case models.QTYPE_NS:
		offset, err = handleNSRecord(response, offset, question)
	case models.QTYPE_PTR:
		offset, err = handlePTRRecord(response, offset, question)
	case models.QTYPE_SOA:
		offset, err = handleSOARecord(response, offset, question)
	case models.QTYPE_TXT:
		offset, err = handleTXTRecord(response, offset, question)
	default:
		err = ErrUnsupportedQueryType
	}

	if err != nil {
		header.RCode = getRCodeForError(err)
		writeHeader(response, header)
		return response[:offset]
	}

	// NOT IMPLEMENTED
	if header.NSCount > 0 {
		offset = writeNSRecords(response, offset, question)
	}

	// NOT IMPLEMENTED
	if header.ARCount > 0 {
		offset = writeAdditionalRecords(response, offset, question)
	}

	// Doesn't matter the length, we will send the response till offset
	return response[:offset]
}

func writeHeader(response []byte, header models.DNSHeader) int {
	offset := 0
	binary.BigEndian.PutUint16(response[offset:], header.PacketID)
	offset += 2

	flags := uint16(1)<<15 | uint16(header.Opcode)<<11 | uint16(header.AA)<<10 |
		uint16(header.TC)<<9 | uint16(header.RD)<<8 |
		uint16(header.RA)<<7 | uint16(header.Z)<<4 | uint16(header.RCode)
	binary.BigEndian.PutUint16(response[offset:], flags)
	offset += 2

	binary.BigEndian.PutUint16(response[offset:], header.QDCount)
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], header.ANCount)
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], header.NSCount)
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], header.ARCount)
	offset += 2

	return offset
}

func writeQuestion(response []byte, offset int, question models.DNSQuestion) int {
	for _, label := range question.QName {
		response[offset] = byte(len(label))
		offset++
		copy(response[offset:], label)
		offset += len(label)
	}
	response[offset] = 0
	offset++

	binary.BigEndian.PutUint16(response[offset:], uint16(question.QType))
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], question.QClass)
	offset += 2

	return offset
}

func handleARecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	qname := joinLabels(question.QName)
	ttl, rdLength, ipAddr := responsehandlers.AHandler(qname)
	if ipAddr == 0 {
		return offset, ErrRecordNotFound
	}

	offset = writeResourceRecord(response, offset, question, models.QTYPE_A, ttl, rdLength)
	binary.BigEndian.PutUint32(response[offset:], ipAddr)
	offset += IPv4AddressLength
	incrementAnswerCount(response)

	return offset, nil
}

func handleAAAARecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	qname := joinLabels(question.QName)
	ttl, rdLength, ipAddr := responsehandlers.AAAA_handler(qname)

	if ipAddr == nil {
		return offset, ErrRecordNotFound
	}

	offset = writeResourceRecord(response, offset, question, models.QTYPE_AAAA, ttl, rdLength)
	copy(response[offset:], ipAddr)
	offset += len(ipAddr)

	incrementAnswerCount(response)
	return offset, nil
}

func handleCNAMERecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later (tough one ig)
	return offset, ErrNotImplemented
}

func handleMXRecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later
	return offset, ErrNotImplemented
}

func handleNSRecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later :D
	return offset, ErrNotImplemented
}

func handlePTRRecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later :D
	return offset, ErrNotImplemented
}

func handleSOARecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later
	return offset, ErrNotImplemented
}

func handleTXTRecord(response []byte, offset int, question models.DNSQuestion) (int, error) {
	// Later
	return offset, ErrNotImplemented
}

// IMPORTANT TO IMPLEMENT BEFORE MAKING IT FUNCTIONAL
func writeNSRecords(response []byte, offset int, question models.DNSQuestion) int {
	// Implement NS records writing
	return offset
}

func writeAdditionalRecords(response []byte, offset int, question models.DNSQuestion) int {
	// Implement additional records writing
	return offset
}
