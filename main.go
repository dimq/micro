package main

import (
	"io/ioutil"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/namsral/flag"
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

func init() {
	flag.CommandLine = flag.NewFlagSetWithEnvPrefix("kafka-set", "KAFKA", flag.ExitOnError)

	flag.StringVar(&brokers, "brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&topics, "topics", "", "Kafka topics to be consumed, as a comma seperated list")
	flag.StringVar(&username, "username", "", "Kafka username to use")
	flag.StringVar(&password, "password", "", "Kafka password to use")
	flag.StringVar(&group, "group", "", "Kafka consumer group definition")
	flag.StringVar(&version, "version", "2.1.1", "Kafka cluster version")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
}

func main() {
	//Setup logger
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//Init database connection
	initDB()

	//Setup consumer for Kafka
	Consume()

}
