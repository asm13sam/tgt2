package main

import (
	"encoding/json"
	"os"
	"strings"
)

var typesTemplates = map[string]string{
	"int":   "\t\"%s\": %s,\n",
	"str":   "\t\"%s\": \"%s\",\n",
	"bool":  "\t\"%s\": %s,\n",
	"float": "\t\"%s\": %s,\n",
}

var operators = map[string]string{
	"eq":  "=",
	"lt":  "<",
	"gt":  ">",
	"lte": "<=",
	"gte": ">=",
	"ne":  "<>",
}

const (
	FormHide = iota
	FormVisible
	FormRequired
)

const (
	MessageHide = iota
	MessageVisible
)

type Model struct {
	Documents     []string             `json:"documents"`
	DocTableItems []string             `json:"doc_table_items"`
	Models        map[string]ItemModel `json:"models"`
}

type FKey struct {
	Table string `json:"table"`
	Field string `json:"field"`
}

type ItemColumn struct {
	Def  interface{} `json:"def"`
	Hum  string      `json:"hum"`
	Form int         `json:"form"`
	Type string      `json:"type"`
}

type ItemModel struct {
	Hum       string                `json:"hum"`
	Rights    string                `json:"rights"`
	Message   int                   `json:"message"`
	FKeys     map[string][]string   `json:"fkeys"`
	Columns   []string              `json:"columns"`
	WColumns  []string              `json:"w_columns"`
	Related   []Related             `json:"related"`
	Registers []Register            `json:"registers"`
	Model     map[string]ItemColumn `json:"model"`
	WModel    map[string]ItemColumn `json:"w_model"`
}

type Related struct {
	Table       string `json:"table"`
	Filter      string `json:"filter"`
	FilterValue string `json:"filter_value"`
}

type Register struct {
	RegField string   `json:"reg_field"`
	ValField []string `json:"val_field"`
	Func     string   `json:"func"`
	Type     string   `json:"type"`
}

var Md Model

func readModel() error {
	data, err := os.ReadFile("models.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	err = decoder.Decode(&Md)
	if err != nil {
		return err
	}
	return nil
}
