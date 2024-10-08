package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
)

func CreateDomain(c *gin.Context) {
	type CreateDomainRequest struct {
		DomainName string `json:"domain_name" binding:"required"`
		ParentID   *uint  `json:"parent_id"`
		IsActive   bool   `json:"is_active"`
		IP         string `json:"ip" binding:"required"`          // IP address for the @ A record
		AdminEmail string `json:"admin_email" binding:"required"` // Admin email for SOA record
		Serial     int    `json:"serial" binding:"required"`      // Serial number for SOA record
		Refresh    int    `json:"refresh" binding:"required"`     // Refresh time for SOA record
		Retry      int    `json:"retry" binding:"required"`       // Retry time for SOA record
		Expire     int    `json:"expire" binding:"required"`      // Expiration time for SOA record
		TTL        int    `json:"ttl" binding:"required"`         // TTL for DNS records
	}

	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain := models.Domain{
		DomainName: req.DomainName,
		ParentID:   req.ParentID,
		IsActive:   req.IsActive,
	}

	if err := database.DB.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create domain"})
		return
	}
}

// // Only IN CLass Allowed ( I don't know how to handle rest as of now )
func CreateRecord(c *gin.Context) {
	type CreateRecordRequest struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		RecordType  string `json:"record_type" binding:"required"`
		RecordName  string `json:"record_name" binding:"required"`
		RecordValue string `json:"record_value" binding:"required"`
		TTL         int    `json:"ttl" binding:"required"`
		Priority    *int   `json:"priority"`
	}

	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

}

// 	var existingRecord models.DNSRecords
// 	if err := database.DB.Where("name = ? AND data = ?", body.Name, body.Data).First(&existingRecord).Error; err == nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": "Conflicting Records"})
// 		return
// 	}

// 	dnsRecord := models.DNSRecords{
// 		Name:    body.Name,
// 		Type:    body.Type,
// 		TTL:     body.TTL,
// 		Data:    body.Data,
// 		BaseURL: body.Base,
// 		Class:   "IN",
// 	}

// 	result := database.DB.Create(&dnsRecord)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Record Successfully Created",
// 		"record":  dnsRecord,
// 	})
// }

func UpdateRecordByID(c *gin.Context) {
	domainId := c.Param("domainId")

	var body struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		RecordType  string `json:"record_type" binding:"required"`
		RecordName  string `json:"record_name" binding:"required"`
		RecordValue string `json:"record_value" binding:"required"`
		TTL         int    `json:"ttl" binding:"required"`
		Priority    *int   `json:"priority"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if domainId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DomainID Parameter is required"})
		return
	}

	var dnsRecord models.DNSRecord

	if err := database.DB.First(&dnsRecord, "domain_id=?", domainId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if body.RecordType != "" {
		dnsRecord.RecordType = body.RecordType
	}
	if body.RecordType != "" {
		dnsRecord.RecordType = body.RecordType
	}
	if body.RecordValue != "" {
		dnsRecord.RecordValue = body.RecordValue
	}
	if body.TTL != 0 {
		dnsRecord.TTL = body.TTL
	}
	if body.Priority != nil {
		dnsRecord.Priority = body.Priority
	}

	if err := database.DB.Save(&dnsRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record successfully updated",
		"record":  dnsRecord,
	})
}

func DeleteDomainByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Parameter is required"})
		return
	}

	var domain models.Domain
	if err := database.DB.First(&domain, "id=?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := database.DB.Delete(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Domain successfully deleted",
	})
}
