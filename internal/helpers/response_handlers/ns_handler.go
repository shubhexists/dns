package responsehandlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shubhexists/dns/cache"
	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
)

func NSHandler(Qname string) (uint32, uint16, []byte) {
	diceDB := cache.NewAPIClient()

	res, err := diceDB.Get(Qname)
	if err != nil {
		fmt.Println("Error", err)
		return 0, 0, nil
	}

	if res == nil {
		var domain models.Domain
		if err := database.DB.Where("domain_name = ?", Qname).First(&domain).Error; err != nil {
			fmt.Println("Domain not found:", err)
			return 0, 0, nil
		}

		var dnsRecord models.DNSRecord
		if err := database.DB.Where("domain_id = ? AND record_type = ?", domain.ID, "NS").First(&dnsRecord).Error; err != nil {
			fmt.Println("A record not found:", err)
			return 0, 0, nil
		}

		cacheData := data{
			Key:   Qname,
			Value: dnsRecord.RecordValue,
			TTL:   strconv.Itoa(dnsRecord.TTL),
		}

		cacheBytes, error := json.Marshal(cacheData)
		if error != nil {
			fmt.Println("Unable to marshal cache data: ", error)
			return 0, 0, nil
		}

		diceDB.Set(Qname, string(cacheBytes))

		res = cacheBytes
	}

	var resData data

	if err := json.Unmarshal(res, &resData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, 0, nil
	}

	ttl, err := strconv.ParseUint(resData.TTL, 10, 32)
	if err != nil {
		fmt.Println("Error converting to uint32:", err)
		return 0, 0, nil
	}

	name_server := strings.Split(resData.Value, ".")

	var rdata []byte

	for _, s := range name_server {
		rdata = append(rdata, byte(len(s)))
		rdata = append(rdata, []byte(s)...)
	}

	rdata = append(rdata, 0x00)

	rdlength := uint16(len(rdata))

	return uint32(ttl), rdlength, rdata
}
