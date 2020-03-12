package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func InitDB() (*gorm.DB, error) {
	var (
		dbHost string
		dbUser string
		dbPass string
		dbName string
		dbPort string
	)

	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USERNAME")
	dbPass = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbPort = os.Getenv("DB_PORT")

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPass)

	db, err := gorm.Open(
		"postgres",
		dbUri)

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&HashTable{})
	return db, nil
}
