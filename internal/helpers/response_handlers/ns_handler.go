package responsehandlers

import (
	"strings"

	"github.com/shubhexists/dns/cache"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
)

func NSHandler(Qname string) (uint32, uint16, []byte) {
	diceDB := cache.NewAPIClient()

	res, err := diceDB.Get(Qname)
	if err != nil {
		Log.Errorln("Error getting domain values:", err)
		return 0, 0, nil
	}

	var nsRecord cache.RecordData
	var exists bool

	if res != nil {
		nsRecord, exists = res["NS"]
	}

	if !exists {
		var domain models.Domain
		if err := database.DB.Where("domain_name = ?", Qname).First(&domain).Error; err != nil {
			Log.Errorln("Domain not found:", err)
			return 0, 0, nil
		}

		var dnsRecord models.DNSRecord
		// TODO: ADD IN README THAT WE WOULD SUPPORT ONLY 1 NS for each DEPLOYMENT
		if err := database.DB.Where("domain_id = ? AND record_type = ?", domain.ID, "NS").First(&dnsRecord).Error; err != nil {
			Log.Errorln("NS record not found:", err)
			return 0, 0, nil
		}

		nsRecord = cache.RecordData{
			Value: dnsRecord.RecordValue,
			TTL:   dnsRecord.TTL,
		}

		cacheData := map[string]cache.RecordData{
			"NS": nsRecord,
		}

		err = diceDB.Set(Qname, cacheData)
		if err != nil {
			Log.Errorln("Unable to set cache data:", err)
			return 0, 0, nil
		}
	}

	name_server := strings.Split(nsRecord.Value, ".")

	var rdata []byte

	for _, s := range name_server {
		rdata = append(rdata, byte(len(s)))
		rdata = append(rdata, []byte(s)...)
	}

	rdata = append(rdata, 0x00)

	rdlength := uint16(len(rdata))

	return uint32(nsRecord.TTL), rdlength, rdata
}
