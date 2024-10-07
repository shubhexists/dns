package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhexists/dns/database"
	"github.com/shubhexists/dns/models"
)

// Only IN CLass Allowed ( I don't know how to handle rest as of now )
func CreateRecord(c *gin.Context) {
	var body struct {
		Base string `json:"base" binding:"required"`
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
		TTL  uint32 `json:"ttl" binding:"required"`
		Data string `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingRecord models.DNSRecords
	if err := database.DB.Where("name = ? AND data = ?", body.Name, body.Data).First(&existingRecord).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Conflicting Records"})
		return
	}

	dnsRecord := models.DNSRecords{
		Name:    body.Name,
		Type:    body.Type,
		TTL:     body.TTL,
		Data:    body.Data,
		BaseURL: body.Base,
		Class:   "IN",
	}

	result := database.DB.Create(&dnsRecord)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record Successfully Created",
		"record":  dnsRecord,
	})
}

func UpdateRecordByID(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Base string  `json:"base"`
		Name string  `json:"name"`
		Type string  `json:"type"`
		TTL  *uint32 `json:"ttl"`
		Data string  `json:"data"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Parameter is required"})
		return
	}

	var dnsRecord models.DNSRecords

	if err := database.DB.First(&dnsRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if body.Name != "" {
		dnsRecord.Name = body.Name
	}
	if body.Type != "" {
		dnsRecord.Type = body.Type
	}
	if body.TTL != nil {
		dnsRecord.TTL = *body.TTL
	}
	if body.Data != "" {
		dnsRecord.Data = body.Data
	}
	if body.Base != "" {
		dnsRecord.BaseURL = body.Base
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

func DeleteRecordByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Parameter is required"})
		return
	}

	var dnsRecord models.DNSRecords
	if err := database.DB.First(&dnsRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := database.DB.Delete(&dnsRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record successfully deleted",
	})
}

func DeleteRecordsByName(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is Required"})
		return
	}

	if err := database.DB.Where("Name = ?", name).Delete(&models.DNSRecords{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Records successfully deleted",
	})
}

func GetRecordByID(c *gin.Context) {
	id := c.Param("id")

	var dnsRecord models.DNSRecords
	if err := database.DB.First(&dnsRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"record": dnsRecord,
	})
}

func GetRecordsByName(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	var dnsRecords []models.DNSRecords
	if err := database.DB.Where("name = ?", name).Find(&dnsRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": dnsRecords,
	})
}
