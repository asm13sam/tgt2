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
	mode := r.URL.Query().Get("mode")
	var res string
	switch mode {
	case "raw":
		res, err = GetItemRaw(params["table"], id)
	default:
		res, err = GetItem(params["table"], id)
	}
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Record not found"))
		return
	}
	Respond(w, http.StatusOK, res)
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res, err := GetItems(params["table"])
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Records not found"))
		return
	}
	Respond(w, http.StatusOK, res)
}

func getFilterHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res, err := GetFilter(params["table"], params["column"], params["operator"], params["value"])
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Record not found"))
		return
	}
	Respond(w, http.StatusOK, res)
}
