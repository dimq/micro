package main

import (
	"fmt"
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
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(&HashTable{})

}

func CheckHashExist(hash string) bool {
	fmt.Println(hash)
	var tempHash HashTable
	if db.Where("hash = ?", hash).First(&tempHash).RecordNotFound() {
		return true
	}
	//if !db.First(&tempHash, "hash = ?", hash).RecordNotFound() {
	//	return true
	//}
	return false
}

func InsertHash(hash string) bool {
	if err := db.Create(&HashTable{Hash: hash}).GetErrors(); err != nil {
		return false
	}
	return true
}
