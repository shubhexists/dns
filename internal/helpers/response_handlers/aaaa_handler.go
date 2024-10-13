package responsehandlers

import (
	"net"

	"github.com/shubhexists/dns/cache"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
)

func AAAA_handler(QName string) (uint32, uint16, []byte) {
	diceDB := cache.NewAPIClient()

	if QName == "" {
		return 0, 0, nil
	}

	res, err := diceDB.Get(QName)
	if err != nil {
		Log.Errorln("Error getting domain values:", err)
		return 0, 0, nil
	}

	var aaaaRecord cache.RecordData
	var exists bool

	if res != nil {
		aaaaRecord, exists = res["AAAA"]
	}

	if !exists {
		var domain models.Domain
		if err := database.DB.Where("domain_name = ?", QName).First(&domain).Error; err != nil {
			Log.Errorln("Domain record not found:", err)
			return 0, 0, nil
		}

		var record models.DNSRecord
		if err := database.DB.Where("domain_id = ? AND record_type = ?", domain.ID, "AAAA").First(&record).Error; err != nil {
			Log.Errorln("AAAA record not found:", err)
			return 0, 0, nil
		}

		aaaaRecord = cache.RecordData{
			Value: record.RecordValue,
			TTL:   record.TTL,
		}

		cacheData := map[string]cache.RecordData{
			"AAAA": aaaaRecord,
		}

		err = diceDB.Set(QName, cacheData)
		if err != nil {
			Log.Errorln("Unable to set cache data:", err)
			return 0, 0, nil
		}
	}

	ip := net.ParseIP(aaaaRecord.Value)
	if ip == nil {
		Log.Errorln("Invalid IP address")
		return 0, 0, nil
	}

	ip = ip.To16()
	if ip == nil || ip.To4() != nil {
		Log.Errorln("Not a valid IPv6 address")
		return 0, 0, nil
	}

	return uint32(aaaaRecord.TTL), 0x001C, ip
}
