package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/database"
	. "github.com/shubhexists/dns/internal/logger"
	"github.com/shubhexists/dns/models"
	"gorm.io/gorm"
)

func CreateDomain(c *gin.Context) {
	Log.Println("Creating domain...")
	var req struct {
		DomainName string `json:"domain_name" binding:"required"`
		ParentID   *uint  `json:"parent_id"`
		IP         string `json:"ip" binding:"required"`
		TTL        int    `json:"ttl" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Log.Errorf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain := models.Domain{
		DomainName: req.DomainName,
		ParentID:   req.ParentID,
		IsActive:   true,
	}

	if err := database.DB.Create(&domain).Error; err != nil {
		Log.Errorf("Could not create domain: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create domain"})
		return
	}

	records := []models.DNSRecord{
		{
			DomainID:    domain.ID,
			RecordType:  "A",
			RecordName:  "@",
			TTL:         req.TTL,
			RecordValue: req.IP,
		},
		{
			DomainID:    domain.ID,
			RecordType:  "CNAME",
			RecordName:  "www",
			TTL:         req.TTL,
			RecordValue: domain.DomainName,
		},
		{
			DomainID:   domain.ID,
			RecordType: "NS",
			RecordName: "@",
			TTL:        3600,
			// TODO: Replace when deployed
			RecordValue: "ns1.shubh.sh",
		},
	}

	if err := database.DB.Create(&records).Error; err != nil {
		Log.Errorf("Could not create DNS records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create DNS records"})
		return
	}

	Log.Println("Domain successfully created")
	c.JSON(http.StatusOK, gin.H{"message": "Domain successfully created"})
}

func CreateRecord(c *gin.Context) {
	Log.Println("Creating DNS record...")
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		RecordType  string `json:"record_type" binding:"required"`
		RecordName  string `json:"record_name" binding:"required"`
		RecordValue string `json:"record_value" binding:"required"`
		TTL         int    `json:"ttl" binding:"required"`
		Priority    *int   `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Log.Errorf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var domain models.Domain
	if err := database.DB.First(&domain, req.DomainID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Log.Errorln("Domain not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		} else {
			Log.Errorf("Error checking domain existence: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking domain existence"})
		}
		return
	}

	newRecord := models.DNSRecord{
		DomainID:    req.DomainID,
		RecordType:  req.RecordType,
		RecordName:  req.RecordName,
		RecordValue: req.RecordValue,
		TTL:         req.TTL,
		Priority:    req.Priority,
	}

	if err := database.DB.Create(&newRecord).Error; err != nil {
		Log.Errorf("Could not create DNS record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create DNS Record"})
		return
	}

	Log.Println("Record created successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Record created successfully", "record": newRecord})
}

func UpdateRecord(c *gin.Context) {
	Log.Println("Updating DNS record...")
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		RecordType  string `json:"record_type" binding:"required"`
		RecordName  string `json:"record_name" binding:"required"`
		RecordValue string `json:"record_value"`
		TTL         *int   `json:"ttl"`
		Priority    *int   `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Log.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dnsRecord models.DNSRecord
	// TODO: Add in README that we don't support multiple values for same record_type and record_name
	if err := database.DB.Where("domain_id = ? AND record_type = ? AND record_name = ?",
		req.DomainID, req.RecordType, req.RecordName).First(&dnsRecord).Error; err != nil {
		Log.Errorln("Record not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.RecordValue != "" {
		updates["record_value"] = req.RecordValue
	}
	if req.TTL != nil {
		updates["ttl"] = *req.TTL
	}
	if req.Priority != nil {
		updates["priority"] = req.Priority
	}

	if len(updates) == 0 {
		Log.Errorln("No valid updates provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid updates provided"})
		return
	}

	if err := database.DB.Model(&dnsRecord).Updates(updates).Error; err != nil {
		Log.Errorf("Failed to update record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	Log.Println("Record successfully updated")
	c.JSON(http.StatusOK, gin.H{
		"message": "Record successfully updated",
		"record":  dnsRecord,
	})
}

func DeleteDomainByID(c *gin.Context) {
	Log.Println("Deleting domain by ID...")
	id := c.Param("id")

	domainID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		Log.Errorf("Invalid domain ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	tx := database.DB.Begin()

	if err := tx.Delete(&models.Domain{}, domainID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			Log.Errorln("Domain not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		} else {
			Log.Errorln("Failed to delete domain: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		}
		return
	}

	result := tx.Where("domain_id = ?", domainID).Delete(&models.DNSRecord{})
	if result.Error != nil {
		tx.Rollback()
		Log.Errorln("Failed to delete DNS records: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS records"})
		return
	}
	Log.Printf("Deleted %d DNS records\n", result.RowsAffected)

	result = tx.Where("domain_id = ?", domainID).Delete(&models.SOARecord{})
	if result.Error != nil {
		tx.Rollback()
		Log.Errorf("Failed to delete SOA record: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SOA record"})
		return
	}
	Log.Printf("Deleted %d SOA records\n", result.RowsAffected)

	if err := tx.Commit().Error; err != nil {
		Log.Errorf("Failed to commit changes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	Log.Println("Domain and all related records successfully deleted")
	c.JSON(http.StatusOK, gin.H{"message": "Domain and all related records successfully deleted"})
}

func GetRecordsByDomainID(c *gin.Context) {
	Log.Println("Retrieving records for domain ID...")
	domainIDParam := c.Param("domain_id")

	domainID, err := strconv.ParseUint(domainIDParam, 10, 32)
	if err != nil {
		Log.Errorf("Invalid domain ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain ID"})
		return
	}

	var records []models.DNSRecord
	if err := database.DB.Where("domain_id = ?", domainID).Find(&records).Error; err != nil {
		Log.Errorf("Error retrieving DNS records: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving DNS records"})
		return
	}

	if len(records) == 0 {
		Log.Println("No records found for this domain")
		c.JSON(http.StatusNotFound, gin.H{"message": "No records found for this domain"})
		return
	}

	Log.Printf("Retrieved %d records for domain ID %d\n", len(records), domainID)
	c.JSON(http.StatusOK, gin.H{"records": records})
}

func DeleteRecord(c *gin.Context) {
	Log.Println("Deleting DNS record...")
	var req struct {
		DomainID   uint   `json:"domain_id" binding:"required"`
		RecordType string `json:"record_type" binding:"required"`
		RecordName string `json:"record_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Log.Errorf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("domain_id = ? AND record_type = ? AND record_name = ?",
		req.DomainID, req.RecordType, req.RecordName).Delete(&models.DNSRecord{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Log.Errorln("Record not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		Log.Errorf("Failed to delete record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	Log.Println("Record successfully deleted")
	c.JSON(http.StatusOK, gin.H{"message": "Record successfully deleted"})
}
