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
	active := r.URL.Query().Get("active")
	items := &Items{name: tableName, mode: mode, active: active}
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
	active := r.URL.Query().Get("active")
	if operator == "like" {
		value = "%" + strings.ToLower(value) + "%"
	}
	items := &FilteredItems{
		Items:        Items{name: tableName, mode: mode, active: active},
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

func getFilterSumHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	filterColumn := params["filter_column"]
	sumColumn := params["sum_column"]
	operator := params["operator"]
	value := params["value"]
	fmt.Println("get sum filter", tableName, sumColumn, filterColumn, operator, value)
	sum, err := GetSumFilter(tableName, sumColumn, filterColumn, operator, value)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, fmt.Sprintf("{\"sum\": %f}", sum))
}

func getBetweenHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	filterColumn := params["column"]
	value_start := params["value_start"]
	value_end := params["value_end"]
	mode := r.URL.Query().Get("mode")
	active := r.URL.Query().Get("active")
	items := &BetweenItems{
		Items:        Items{name: tableName, mode: mode, active: active},
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

func getSumHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tableName := params["table"]
	filterColumn := params["column"]
	fmt.Println("get sum", tableName, filterColumn)
	sum, err := GetSum(tableName, filterColumn)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Can't update item"))
		return
	}
	Respond(w, http.StatusOK, fmt.Sprintf("{\"sum\": %f}", sum))
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
