package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
)

// We will NOT have a domain ID when calling this
func CreateDomain(c *gin.Context) {
	type CreateDomainRequest struct {
		DomainName string `json:"domain_name" binding:"required"`
		ParentID   *uint  `json:"parent_id"`
		IP         string `json:"ip" binding:"required"`          // IP address for the @ A record
		AdminEmail string `json:"admin_email" binding:"required"` // Admin email for SOA record
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
		IsActive:   true,
	}

	if err := database.DB.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create domain"})
		return
	}

	// Making a root record for the new domain
	rootRecord := models.DNSRecord{
		DomainID:    domain.ID,
		RecordType:  "A",
		RecordName:  "@",
		TTL:         req.TTL,
		RecordValue: req.IP,
		Priority:    nil,
	}

	if err := database.DB.Create(&rootRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create A record for root domain"})
		return
	}

	// CNAME www to the root domain name
	wwwRecord := models.DNSRecord{
		DomainID:   domain.ID,
		RecordType: "CNAME",
		RecordName: "www",
		TTL:        req.TTL,
		// Confirm this
		RecordValue: domain.DomainName,
		Priority:    nil,
	}

	if err := database.DB.Create(&wwwRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cname for WWW"})
		return
	}

	// Create a NS Record and SOA record

	soaRecord := models.SOARecord{
		DomainID: domain.ID,
		// Replace this when deplyed or take from env
		PrimaryNS: "",
		// shubh622005@gmail.com should be written as shubh622005.gmail.com
		// NO "@"
		AdminEmail: "",
		TTL:        req.TTL,
	}

	if err := database.DB.Create(&soaRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create SOA record"})
		return
	}

	// We will only support one NS server

	nsServer := models.Nameserver{
		DomainID: domain.ID,
		// Replace this when deplyed or take from env
		NSName: "",
	}

	if err := database.DB.Create(&nsServer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create NS Record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Domain successfully created"})
}

// We will have a DomainID for sure while calling this
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	newRecord := models.DNSRecord{
		DomainID:   req.DomainID,
		RecordType: req.RecordType,
		TTL:        req.TTL,
		Priority:   req.Priority,
	}

	if err := database.DB.Create(&newRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create DNS Record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Records created successfully"})
}

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
