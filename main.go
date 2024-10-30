package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

const (
	INT = iota
	STRING
	BOOL
	FLOAT
)

var (
	db               *sql.DB
	existingTables   []string
	tablesColumnsNum map[string]int
	tablesColumns    map[string][]TableColumn
	tablesColumns1   map[string](map[string]TableColumn)
)

var typesTemplates = []string{
	"\t\"%s\": %s,\n",
	"\t\"%s\": \"%s\",\n",
	"\t\"%s\": %s,\n",
	"\t\"%s\": %s,\n",
}

var operands = map[string]string{
	"eq":  "=",
	"lt":  "<",
	"gt":  ">",
	"lte": "<=",
	"gte": ">=",
}

type TableColumn struct {
	cname string
	ctype int
}

func makeExistingTables() {
	existingTables = []string{"measure", "matherial_group", "matherial"}
	tablesColumnsNum = map[string]int{"measure": 4, "matherial_group": 5, "matherial": 12}
	tablesColumns = map[string][]TableColumn{
		"measure": {TableColumn{"id", INT}, TableColumn{"name", STRING},
			TableColumn{"full_name", STRING}, TableColumn{"is_active", BOOL}},
		"matherial_group": {TableColumn{"id", INT}, TableColumn{"name", STRING},
			TableColumn{"matherial_group_id", INT}, TableColumn{"is_active", BOOL}},
		"matherial": {TableColumn{"id", INT}, TableColumn{"name", STRING}, TableColumn{"full_name", STRING},
			TableColumn{"matherial_group_id", INT}, TableColumn{"measure_id", INT},
			TableColumn{"color_group_id", INT}, TableColumn{"price", FLOAT},
			TableColumn{"cost", FLOAT}, TableColumn{"total", FLOAT},
			TableColumn{"barcode", STRING}, TableColumn{"count_type_id", INT}, TableColumn{"is_active", BOOL}},
	}
	tablesColumns1 = map[string](map[string]TableColumn){
		"measure": {
			"id":        TableColumn{"id", INT},
			"name":      TableColumn{"name", STRING},
			"full_name": TableColumn{"full_name", STRING},
			"is_active": TableColumn{"is_active", BOOL},
		},
		"matherial_group": {
			"id":                 TableColumn{"id", INT},
			"name":               TableColumn{"name", STRING},
			"matherial_group_id": TableColumn{"matherial_group_id", INT},
			"is_active":          TableColumn{"is_active", BOOL},
		},
		"matherial": {
			"id":                 TableColumn{"id", INT},
			"name":               TableColumn{"name", STRING},
			"full_name":          TableColumn{"full_name", STRING},
			"matherial_group_id": TableColumn{"matherial_group_id", INT},
			"measure_id":         TableColumn{"measure_id", INT},
			"color_group_id":     TableColumn{"color_group_id", INT},
			"price":              TableColumn{"price", FLOAT},
			"cost":               TableColumn{"cost", FLOAT},
			"total":              TableColumn{"total", FLOAT},
			"barcode":            TableColumn{"barcode", STRING},
			"count_type_id":      TableColumn{"count_type_id", INT},
			"is_active":          TableColumn{"is_active", BOOL}},
	}
}

func DBconnect(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	return err
}

func DbClose() error {
	return db.Close()
}

func testForExistingTable(tableName string) bool {
	for _, f := range existingTables {
		if tableName == f {
			return true
		}
	}
	return false
}

func makeSql(tableName, operator string, params ...string) (string, error) {
	var res string
	var err error
	if !testForExistingTable(tableName) {
		return res, fmt.Errorf("table %s is not exist", tableName)
	}
	switch operator {
	case "by_id":
		res = fmt.Sprintf("SELECT * FROM %s WHERE id=?", tableName)
	case "all":
		res = fmt.Sprintf("SELECT * FROM %s WHERE is_active=1", tableName)
	case "filter":
		res = fmt.Sprintf("SELECT * FROM %s WHERE is_active=1 AND %s %s ?", tableName, params[0], operands[params[1]])
	default:
		err = fmt.Errorf("operator %s is not exist", operator)
	}
	return res, err
}

func makeHolder(number int) []interface{} {
	m := []interface{}{}
	for i := 0; i < number; i++ {
		s := ""
		m = append(m, &s)
	}
	return m
}

func makeJsonMap(tableName string, holder []interface{}) string {
	js := ""
	for i, s := range holder {
		v := *s.(*string)
		t := tablesColumns[tableName][i].ctype
		if t == BOOL {
			if v == "0" {
				v = "false"
			} else {
				v = "true"
			}
		}
		js += fmt.Sprintf(typesTemplates[t], tablesColumns[tableName][i].cname, v)
	}
	return fmt.Sprintf("{\n%s\n}", js[:len(js)-2])
}

func query(ctype int, value, sel string) (*sql.Rows, error) {
	switch ctype {
	case INT:
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return db.Query(sel, num)
	case STRING:
		return db.Query(sel, value)
	case BOOL:
		return db.Query(sel, value == "true" || value == "1")
	case FLOAT:
		fnum, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return db.Query(sel, fnum)
	default:
		return nil, fmt.Errorf("unknown type code %d for value %s", ctype, value)
	}
}

func GetFilter(tableName, filterColumn, operator, value string) (string, error) {
	sel, err := makeSql(tableName, "filter", filterColumn, operator)
	if err != nil {
		return "", err
	}

	holder := makeHolder(tablesColumnsNum[tableName])
	rows, err := query(tablesColumns1[tableName][filterColumn].ctype, value, sel)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	jsMaps := ""
	for rows.Next() {
		err = rows.Scan(holder...)
		if err != nil {
			return "", err
		}
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(tableName, holder))
	}
	jsList := fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	return jsList, nil
}

func GetAll(tableName string) (string, error) {
	sel, err := makeSql(tableName, "all")
	if err != nil {
		return "", err
	}

	holder := makeHolder(tablesColumnsNum[tableName])
	rows, err := db.Query(sel)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	jsMaps := ""
	for rows.Next() {
		err = rows.Scan(holder...)
		if err != nil {
			return "", err
		}
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(tableName, holder))
	}
	jsList := fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	return jsList, nil
}

func Get(tableName string, id int) (string, error) {
	sel, err := makeSql(tableName, "by_id")
	if err != nil {
		return "", err
	}
	holder := makeHolder(tablesColumnsNum[tableName])

	row := db.QueryRow(sel, id)
	err = row.Scan(holder...)
	if err != nil {
		return "", err
	}
	jsMap := makeJsonMap(tableName, holder)
	return jsMap, nil
}

func Respond(w http.ResponseWriter, code int, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func makeError(err error, s string) string {
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%s: %s\"}", s, err)
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
	res, err := Get(params["table"], id)
	if err != nil {
		Respond(w, http.StatusInternalServerError, makeError(err, "Record not found"))
		return
	}
	Respond(w, http.StatusOK, res)
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res, err := GetAll(params["table"])
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

func main() {
	makeExistingTables()

	fmt.Println("Hi")
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
	fmt.Println("Server is running on port 8777...")
	log.Fatal(http.ListenAndServe(":8777", r))
}
