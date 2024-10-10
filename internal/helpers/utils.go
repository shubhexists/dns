package helpers

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/shubhexists/dns/models"
)

func GenerateSOASerial(revision int) (string, error) {
	if revision < 1 || revision > 99 {
		return "", fmt.Errorf("invalid revision number: must be between 01 and 99")
	}

	currentTime := time.Now()
	datePart := currentTime.Format("20060102")

	revisionPart := fmt.Sprintf("%02d", revision)

	serial := datePart + revisionPart
	return serial, nil
}

func writeResourceRecord(response []byte, offset int, question models.DNSQuestion, qType models.QType, ttl uint32, rdLength uint16) int {
	binary.BigEndian.PutUint16(response[offset:], PointerCompression)
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], uint16(qType))
	offset += 2
	binary.BigEndian.PutUint16(response[offset:], question.QClass)
	offset += 2
	binary.BigEndian.PutUint32(response[offset:], ttl)
	offset += 4
	binary.BigEndian.PutUint16(response[offset:], rdLength)
	offset += 2
	return offset
}

func getRCodeForError(err error) uint8 {
	switch err {
	case ErrRecordNotFound:
		// Record not found
		return 3
	case ErrNotImplemented:
		// Maybe later
		return 4
	default:
		// Any other error I don't know about
		return 2
	}
}

func joinLabels(labels []string) string {
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += "."
		}
		result += label
	}
	return result
}

func incrementAnswerCount(response []byte) {
	answerCount := binary.BigEndian.Uint16(response[6:8])
	binary.BigEndian.PutUint16(response[6:8], answerCount+1)
}
