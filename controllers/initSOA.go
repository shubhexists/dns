package controllers

import (
	"strconv"

	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/internal/helpers"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
)

func CheckForSOA() (int, error) {
	var isSOA models.SOARecord

	sno, err := helpers.GenerateSOASerial(1)

	num, err := strconv.Atoi(sno)

	if err != nil {
		return num, err
	}

	if err != nil {
		Log.Errorln("Failed to generate Serial Number")
		return num, err
	}

	if err := database.DB.Where("serial = ?", sno).First(&isSOA).Error; err != nil {
		Log.Warnln("No SOA Record Found. Creating a new SOA Record for Name Server.")
		return num, err
	}

	return num, nil
}
