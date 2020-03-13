package main

type Message struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Consumer struct {
	ready chan bool
}
