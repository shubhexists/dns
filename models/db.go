package models

import (
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	DomainName string  `gorm:"size:255;not null"`
	ParentID   *uint   `gorm:"default:null;constraint:OnDelete:CASCADE"`
	IsActive   bool    `gorm:"default:true"`
	Parent     *Domain `gorm:"foreignKey:ParentID"`
}

type SOARecord struct {
	gorm.Model
	DomainID   uint   `gorm:"not null;constraint:OnDelete:CASCADE"`
	PrimaryNS  string `gorm:"size:255;not null"`
	AdminEmail string `gorm:"size:255;not null"`
	Serial     int    `gorm:"not null"`
	Refresh    int    `gorm:"default:86400"`
	Retry      int    `gorm:"default:7200"`
	Expire     int    `gorm:"default:3600000"`
	TTL        int    `gorm:"default:86400"`
	Domain     Domain `gorm:"foreignKey:DomainID"`
}

type DNSRecord struct {
	gorm.Model
	DomainID    uint   `gorm:"not null;constraint:OnDelete:CASCADE"`
	RecordType  string `gorm:"size:10;not null"`
	RecordName  string `gorm:"size:255;not null"`
	RecordValue string `gorm:"size:255;not null"`
	TTL         int    `gorm:"default:3600"`
	Priority    *int   `gorm:"default:null"`
	Domain      Domain `gorm:"foreignKey:DomainID"`
}

type Nameserver struct {
	gorm.Model
	DomainID uint   `gorm:"not null;constraint:OnDelete:CASCADE"`
	NSName   string `gorm:"size:255;not null"`
	Domain   Domain `gorm:"foreignKey:DomainID"`
}
