package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func Respond(w http.ResponseWriter, code int, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func strError(err error) string {
	return strings.Replace(err.Error(), "\"", "'", -1)
}

func makeError(err error, s string) string {
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%s: %s\"}", s, strError(err))
	}
	return fmt.Sprintf("{\"error\": \"%s\"}", s)
}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		Respond(w, http.StatusBadRequest, makeError(err, "Invalid ID"))
		return
	}
	tableName := params["table"]
	mode := r.URL.Query().Get("mode")
	item := &Item{name: tableName, id: id, mode: mode}
	err = item.Get()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Record not found"))
		return
	}
	Respond(w, http.StatusOK, item.jsonMap)
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	mode := r.URL.Query().Get("mode")
	items := &Items{name: tableName, mode: mode}
	err := items.Get()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, items.jsonList)
}

func getFilterHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	filterColumn := params["column"]
	operator := params["operator"]
	value := params["value"]
	mode := r.URL.Query().Get("mode")
	items := &FilteredItems{
		Items:        Items{name: tableName, mode: mode},
		filterColumn: filterColumn,
		operator:     operator,
		value:        value,
	}
	err := items.Get()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, items.jsonList)
}
