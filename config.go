package main

import (
	"encoding/json"
	"os"
	"strings"
)

var (
	tablesColumnsRaw map[string][]TableColumn
	// tablesColumnsRawMap map[string](map[string]TableColumn)
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

type TableColumn struct {
	cname string
	ctype int
}

func readConfig() {
	tablesColumnsRaw = map[string][]TableColumn{
		"measure": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"is_active", BOOL},
		},
		"matherial_group": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"position", INT},
			TableColumn{"is_active", BOOL},
		},
		"matherial": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"measure_id", INT},
			TableColumn{"color_group_id", INT},
			TableColumn{"price", FLOAT},
			TableColumn{"cost", FLOAT},
			TableColumn{"total", FLOAT},
			TableColumn{"barcode", STRING},
			TableColumn{"count_type_id", INT},
			TableColumn{"is_active", BOOL}},
	}
	// tablesColumnsRawMap = map[string](map[string]TableColumn){
	// 	"measure": {
	// 		"id":        TableColumn{"id", INT},
	// 		"name":      TableColumn{"name", STRING},
	// 		"full_name": TableColumn{"full_name", STRING},
	// 		"is_active": TableColumn{"is_active", BOOL},
	// 	},
	// 	"matherial_group": {
	// 		"id":                 TableColumn{"id", INT},
	// 		"name":               TableColumn{"name", STRING},
	// 		"matherial_group_id": TableColumn{"matherial_group_id", INT},
	// 		"position":           TableColumn{"position", INT},
	// 		"is_active":          TableColumn{"is_active", BOOL},
	// 	},
	// 	"matherial": {
	// 		"id":                 TableColumn{"id", INT},
	// 		"name":               TableColumn{"name", STRING},
	// 		"full_name":          TableColumn{"full_name", STRING},
	// 		"matherial_group_id": TableColumn{"matherial_group_id", INT},
	// 		"measure_id":         TableColumn{"measure_id", INT},
	// 		"color_group_id":     TableColumn{"color_group_id", INT},
	// 		"price":              TableColumn{"price", FLOAT},
	// 		"cost":               TableColumn{"cost", FLOAT},
	// 		"total":              TableColumn{"total", FLOAT},
	// 		"barcode":            TableColumn{"barcode", STRING},
	// 		"count_type_id":      TableColumn{"count_type_id", INT},
	// 		"is_active":          TableColumn{"is_active", BOOL}},
	// }

}
