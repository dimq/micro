package main

import (
	"io/ioutil"
	"os"
)

var (
	brokers  = ""
	topics   = ""
	username = ""
	password = ""
	group    = ""
	version  = ""
)

func main() {
	//Setup logger
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//Setup consumer for Kafka
	Consume()

}
