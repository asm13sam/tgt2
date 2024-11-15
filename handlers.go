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
	fmt.Println("get item", tableName, id, mode)
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
	fmt.Println("get items", tableName, mode)
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
	fmt.Println("get filter", tableName, filterColumn, operator, value, mode)
	err := items.Get()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, items.jsonList)
}

func getBetweenHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	filterColumn := params["column"]
	value_start := params["value_start"]
	value_end := params["value_end"]
	mode := r.URL.Query().Get("mode")
	items := &BetweenItems{
		Items:        Items{name: tableName, mode: mode},
		filterColumn: filterColumn,
		value_start:  value_start,
		value_end:    value_end,
	}
	fmt.Println("get between", tableName, filterColumn, value_start, value_end, mode)
	err := items.Get()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, items.jsonList)
}

func putItemHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	body := r.Body
	fmt.Println("update item", tableName)
	err := UpdateItem(tableName, body)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Can't update item"))
		return
	}
	Respond(w, http.StatusOK, "{}")
}

func postItemHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	body := r.Body
	fmt.Println("create item", tableName)
	err := CreateItem(tableName, body)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Can't create item"))
		return
	}
	Respond(w, http.StatusOK, "{}")
}

func delItemHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		Respond(w, http.StatusBadRequest, makeError(err, "Invalid ID"))
		return
	}
	tableName := params["table"]
	mode := r.URL.Query().Get("mode")
	fmt.Println("delete item", tableName, id, mode)
	err = DeleteItem(tableName, id, mode)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Record not found"))
		return
	}
	Respond(w, http.StatusOK, "{}")
}

func createBaseHandler(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("create base")
	err := CreateBase(body)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Can't create database"))
		return
	}
	Respond(w, http.StatusOK, "{}")
}

func getModelsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get models")
	models, err := GetModels()
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Can't read models"))
		return
	}
	Respond(w, http.StatusOK, models)
}
