package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhexists/dns/database"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/createRecords", CreateRecord)
	return r
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, db_error := database.ConnectToDB()
	if db_error != nil {
		log.Fatalf("Error connecting to database: %v", db_error)
	}

	code := m.Run()

	db.Exec("TRUNCATE TABLE dns_records")

	os.Exit(code)
}

func TestCreateRecord(t *testing.T) {
	if database.DB == nil {
		t.Fatal("Database connection is not initialized")
	}

	tx := database.DB.Begin()
	defer tx.Rollback()

	r := setupRouter()

	t.Run("Success Record", func(t *testing.T) {
		body := `{"base": "example.com", "name": "test", "type": "A", "ttl": 3600, "data": "192.0.2.1"}`
		req, _ := http.NewRequest(http.MethodPost, "/createRecords", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Record Successfully Created")

		// var record models.DNSRecords
		// err := tx.Where("name = ?", "test").First(&record).Error
		// assert.NoError(t, err)
		// assert.Equal(t, "test", record.Name)
	})

	t.Run("Bad Input", func(t *testing.T) {
		body := `{"base":"example.com"}`
		// bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/createRecords", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert response for bad input
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, "Invalid input", response["error"])
	})

	t.Run("Conflict", func(t *testing.T) {
		body := map[string]interface{}{
			"base": "example.com",
			"name": "www.example.com",
			"type": "A",
			"ttl":  300,
			"data": "192.168.1.1",
		}
		bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/createRecords", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert response for conflict
		assert.Equal(t, http.StatusConflict, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Nil(t, err)
		assert.Equal(t, "Conflicting Records", response["error"])
	})
}

func TestDeleteRecordByID(t *testing.T) {
	if database.DB == nil {
		t.Fatal("Database connection is not initialized")
	}

	tx := database.DB.Begin()
	defer tx.Rollback()

	r := setupRouter()

	t.Run("Success Record Deletion", func(t *testing.T) {
		// Insert a test record
		// testRecord := models.DNSRecords{
		// 	BaseURL: "example.com", Name: "test", Type: "A", TTL: 3600, Data: "192.0.2.1",
		// }
		// tx.Create(&testRecord)

		// req, _ := http.NewRequest(http.MethodDelete, "/deleteRecordByID/"+strconv.FormatUint(uint64(testRecord.ID), 10), nil)

		w := httptest.NewRecorder()
		// r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Record successfully deleted")

		// var deletedRecord models.DNSRecords
		// err := tx.First(&deletedRecord, testRecord.ID).Error
		// assert.Error(t, err) // The record should not be found after deletion
	})

	t.Run("Record Not Found", func(t *testing.T) {
		// Non-existent ID
		req, _ := http.NewRequest(http.MethodDelete, "/deleteRecordByID/999999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Record not found")
	})

	t.Run("Missing ID Parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/deleteRecordByID", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID Parameter is required")
	})
}
