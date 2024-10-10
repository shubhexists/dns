package responsehandlers

import (
	"encoding/json"
	"net"
	"strconv"

	"github.com/shubhexists/dns/cache"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
)

type data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

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

	if res == nil {
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

		cacheData := data{
			Key:   Qname,
			Value: dnsRecord.RecordValue,
			TTL:   strconv.Itoa(dnsRecord.TTL),
		}

		cacheBytes, error := json.Marshal(cacheData)
		if error != nil {
			Log.Errorln("Unable to marshal cache data: ", error)
			return 0, 0, 0
		}

		diceDB.Set(Qname, string(cacheBytes))

		res = cacheBytes
	}

	var resData data

	if err := json.Unmarshal(res, &resData); err != nil {
		Log.Errorln("Error parsing JSON:", err)
		return 0, 0, 0
	}

	ttl, err := strconv.ParseUint(resData.TTL, 10, 32)
	if err != nil {
		Log.Errorln("Error converting to uint32:", err)
		return 0, 0, 0
	}

	ip := net.ParseIP(resData.Value)
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

	return uint32(ttl), 0x0004, ipBytes
}
