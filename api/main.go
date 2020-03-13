package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port string = "8080"
)

func handlerMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Println(string(body))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/message", handlerMessage)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
