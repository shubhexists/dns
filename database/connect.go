package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getDBURLFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", errors.New("DB_HOST environment variable is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return "", errors.New("DB_USER environment variable is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "", errors.New("DB_PASSWORD environment variable is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return "", errors.New("DB_NAME environment variable is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("DB_PORT environment variable is not set")
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		return "", errors.New("DB_SSLMODE environment variable is not set")
	}

	timeZone := os.Getenv("DB_TIMEZONE")
	if timeZone == "" {
		return "", errors.New("DB_TIMEZONE environment variable is not set")
	}

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslMode, timeZone)

	return dbURL, nil
}

func ConnectToDB() (*gorm.DB, error) {
	dsn, err := getDBURLFromEnv()
	if err != nil {
		return nil, err
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully")
	return DB, nil
}
