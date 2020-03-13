package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port string = "8080"
	host string = "localhost"
)

type Message struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func handlerMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		msg, err := ParseMessage(body)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(fmt.Sprintf("%+v\n", msg))
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	parseCLI()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/message", handlerMessage)
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(host+":"+port, router))
}

func ParseMessage(msg []byte) (Message, error) {
	var (
		m                  Message
		invalidFormatError = errors.New("invalid format")
	)

	if err := json.Unmarshal([]byte(msg), &m); err != nil {
		return m, invalidFormatError
	}

	if m.ID == 0 {
		return Message{}, invalidFormatError
	}

	if m.Code == "" {
		return Message{}, invalidFormatError
	}

	if m.Message == "" {
		return Message{}, invalidFormatError
	}

	return m, nil
}
