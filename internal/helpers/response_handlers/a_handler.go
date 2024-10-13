package responsehandlers

import (
	"net"

	"github.com/shubhexists/dns/cache"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
)

func AHandler(Qname string) (uint32, uint16, uint32) {
	diceDB := cache.NewAPIClient()

	if Qname == "" {
		Log.Errorln("Error: Invalid QName Value")
		return 0, 0, 0
	}

	res, err := diceDB.Get(Qname)
	if err != nil {
		Log.Errorln("Error", err)
		return 0, 0, 0
	}

	if res == nil || res["A"].Value == "" {
		var domain models.Domain
		if err := database.DB.Where("domain_name = ?", Qname).First(&domain).Error; err != nil {
			Log.Errorln("Domain not found:", err)
			return 0, 0, 0
		}

		var dnsRecord models.DNSRecord
		if err := database.DB.Where("domain_id = ? AND record_type = ?", domain.ID, "A").First(&dnsRecord).Error; err != nil {
			Log.Errorln("A record not found:", err)
			return 0, 0, 0
		}

		cacheData := map[string]cache.RecordData{
			"A": {
				Value: dnsRecord.RecordValue,
				TTL:   dnsRecord.TTL,
			},
		}

		err = diceDB.Set(Qname, cacheData)
		if err != nil {
			Log.Errorln("Unable to set cache data: ", err)
			return 0, 0, 0
		}

		res = cacheData
	}

	aRecord, exists := res["A"]
	if !exists {
		Log.Errorln("A record not found in cache")
		return 0, 0, 0
	}

	ip := net.ParseIP(aRecord.Value)
	if ip == nil {
		Log.Errorln("Invalid IP address")
		return 0, 0, 0
	}

	ip = ip.To4()
	if ip == nil {
		Log.Errorln("Not a valid IPv4 address")
		return 0, 0, 0
	}

	ipBytes := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])

	return uint32(aRecord.TTL), 0x0004, ipBytes
}
