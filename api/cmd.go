package main

import "github.com/namsral/flag"

func parseCLI() {
	flag.StringVar(&host, "host", "localhost", "Adress to bind to the api")
	flag.StringVar(&port, "port", "8080", "Port to bind for the api")
	flag.Parse()
}
