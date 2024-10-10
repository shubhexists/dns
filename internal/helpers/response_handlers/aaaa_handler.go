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

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

func AAAA_handler(QName string) (uint32, uint16, []byte) {
	diceDB := cache.NewAPIClient()

	if QName == "" {
		return 0, 0, nil
	}

	res, err := diceDB.Get(QName)
	if err != nil {
		Log.Errorln("Error getting domain values")
		return 0, 0, nil
	}

	if res == nil {
		var domain models.Domain
		if err := database.DB.Where("domain_name = ?", QName).First(&domain).Error; err != nil {
			Log.Errorln("Domain record not found:", err)
			return 0, 0, nil
		}

		var record models.DNSRecord
		if err := database.DB.Where("domain_id = ? AND record_type = ?", record.DomainID, "AAAA").First(&record).Error; err != nil {
			Log.Errorln("AAAA record not found:", err)
			return 0, 0, nil
		}

		cacheData := data{
			Key:   QName,
			Value: record.RecordValue,
			TTL:   strconv.Itoa(record.TTL),
		}

		cacheBytes, error := json.Marshal(cacheData)
		if error != nil {
			Log.Errorln("Unable to marshal cache data: ", error)
			return 0, 0, nil
		}

		diceDB.Set(QName, string(cacheBytes))

		res = cacheBytes
	}

	var resData data

	if err := json.Unmarshal(res, &resData); err != nil {
		Log.Errorln("Error parsing JSON:", err)
		return 0, 0, nil
	}

	ttl, err := strconv.ParseUint(resData.TTL, 10, 128)
	if err != nil {
		Log.Errorln("Error converting to uint32:", err)
		return 0, 0, nil
	}

	ip := net.ParseIP(resData.Value)
	if ip == nil {
		Log.Errorln("Invalid IP address")
		return 0, 0, nil
	}

	ip = ip.To16()
	if ip == nil || ip.To4() != nil {
		Log.Errorln("Not a valid IPv6 address")
		return 0, 0, nil
	}
	return uint32(ttl), 0x0006, ip
}
