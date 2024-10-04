package helpers

import "github.com/shubhexists/dns/models"

// NOTE -  Replace the params by header and question. I added _ to remove warnings
func BuildDNSResponse(_ models.DNSHeader, _ models.DNSQuestion) []byte {
	response := make([]byte, 512)

	// Implement response querying and stuff

	return response
}
