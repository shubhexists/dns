package models

import (
	"gorm.io/gorm"
)

type DNSRecords struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Type  string `gorm:"not null"`
	TTL   uint32 `gorm:"not null"`
	Data  string `gorm:"not null"`
	Class string `gorm:"not null"`
}
