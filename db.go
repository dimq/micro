package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() {
	var err error
	db, err = gorm.Open(
		"postgres",
		"host="+os.Getenv("HOST")+" port="+os.Getenv("PORT")+" user="+os.Getenv("USER")+
			" dbname="+os.Getenv("DBNAME")+" sslmode=disable password="+
			os.Getenv("PASSWORD"))

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&HashTable{})

}

func CheckHashExist(hash string) bool {
	var tempHash HashTable
	if db.Where("hash = ?", hash).First(&tempHash).RecordNotFound() {
		return true
	}
	return false
}

func InsertHash(hash string) bool {
	if err := db.Create(&HashTable{Hash: hash}).GetErrors(); len(err) != 0 {
		return false
	}
	return true
}

func CloseDB() {
	db.Close()
}
