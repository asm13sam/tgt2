package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	err := readModel()
	if err != nil {
		log.Fatal(err)
	}

	err = DBconnect("base.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DbClose()

	r := mux.NewRouter()

	// Define the endpoints
	r.HandleFunc("/api/models", getModelsHandler).Methods("GET")
	r.HandleFunc("/api/new_base", createBaseHandler).Methods("POST")

	r.HandleFunc("/api/{table}", getItemsHandler).Methods("GET")
	r.HandleFunc("/api/{table}/{id:[0-9]+}", getItemHandler).Methods("GET")
	r.HandleFunc("/api/{table}", putItemHandler).Methods("PUT")
	r.HandleFunc("/api/{table}", postItemHandler).Methods("POST")
	r.HandleFunc("/api/{table}/{id:[0-9]+}", delItemHandler).Methods("DELETE")
	r.HandleFunc("/api/filter/{table}/{column}/{operator}/{value}", getFilterHandler).Methods("GET")
	r.HandleFunc("/api/between/{table}/{column}/{value_start}/{value_end}", getBetweenHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8877...")
	log.Fatal(http.ListenAndServe(":8877", r))
}
