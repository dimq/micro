package main

import (
	"io/ioutil"
	"os"

	"github.com/dimq/micro/models"
	"github.com/jinzhu/gorm"
)

var (
	db       *gorm.DB
	brokers  string
	topics   string
	username string
	password string
	group    string
	version  string
)

func main() {
	var err error

	parseCLI()

	//Setup logger
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//Init database connection
	db, err = models.InitDB()
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()

	//Setup consumer for Kafka
	Consume()

}
