package database

import (
	"errors"
	"fmt"
	"os"

	. "github.com/shubhexists/dns/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getDBURLFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		Log.Errorln("Env variable DB_HOST doesn't exist")
		return "", errors.New("DB_HOST environment variable is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		Log.Errorln("Env variable DB_USER doesn't exist")
		return "", errors.New("DB_USER environment variable is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		Log.Errorln("Env variable DB_PASSWORD doesn't exist")
		return "", errors.New("DB_PASSWORD environment variable is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		Log.Errorln("Env variable DB_NAME doesn't exist")
		return "", errors.New("DB_NAME environment variable is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		Log.Errorln("Env variable DB_USER doesn't exist")
		return "", errors.New("DB_PORT environment variable is not set")
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		Log.Errorln("Env variable DB_SSLMODE doesn't exist")
		return "", errors.New("DB_SSLMODE environment variable is not set")
	}

	timeZone := os.Getenv("DB_TIMEZONE")
	if timeZone == "" {
		Log.Errorln("Env variable DB_TIMEZONE doesn't exist")
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

	Log.Println("Connected to the database successfully")
	return DB, nil
}
