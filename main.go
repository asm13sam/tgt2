package main

import (
	"log"
	"net/http"
	"time"

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
	r.HandleFunc("/api/sum/{table}/{column}", getSumHandler).Methods("GET")
	r.HandleFunc("/api/sum/filter/{table}/{sum_column}/{filter_column}/{operator}/{value}", getFilterSumHandler).Methods("GET")
	r.HandleFunc("/api/filter/{table}/{column}/{operator}/{value}", getFilterHandler).Methods("GET")
	r.HandleFunc("/api/between/{table}/{column}/{value_start}/{value_end}", getBetweenHandler).Methods("GET")

	// Start the server

	server := http.Server{
		Addr:         ":8877",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}

	err = server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

	// fmt.Println("Server is running on port 8877...")
	// log.Fatal(http.ListenAndServe(":8877", r))
}
