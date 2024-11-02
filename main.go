package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	readConfig()

	err := DBconnect("base.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DbClose()

	r := mux.NewRouter()

	// Define the endpoints
	r.HandleFunc("/api/{table}", getItemsHandler).Methods("GET")
	r.HandleFunc("/api/{table}/{id:[0-9]+}", getItemHandler).Methods("GET")
	r.HandleFunc("/api/filter/{table}/{column}/{operator}/{value}", getFilterHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8877...")
	log.Fatal(http.ListenAndServe(":8877", r))
}
