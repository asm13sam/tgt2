package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	name    string
	id      int
	mode    string
	query   string
	holder  []interface{}
	jsonMap string
}

func (i *Item) Get() error {
	if err := i.makeQuery(); err != nil {
		return err
	}
	row := db.QueryRow(i.query, i.id)
	if err := row.Scan(i.holder...); err != nil {
		return err
	}
	i.jsonMap = makeJsonMap(i.name, i.holder, i.mode)
	return nil
}

func (i *Item) makeQuery() error {
	if err := testForExistingTable(i.name); err != nil {
		return err
	}
	var addNum int
	var sels, joins string
	if i.mode != "raw" {
		sels, joins, addNum = makeAddJoins(i.name)
	}
	i.query = fmt.Sprintf("SELECT %s.* %s FROM %s %s WHERE %s.id=?",
		i.name, sels, i.name, joins, i.name)
	i.holder = makeHolder(len(Md.Models[i.name].Columns) + addNum)
	return nil
}

type Items struct {
	name     string
	mode     string
	active   string
	query    string
	holder   []interface{}
	jsonList string
}

func (i *Items) Get() error {
	if err := i.makeQuery(); err != nil {
		return err
	}
	rows, err := i.getRows()
	if err != nil {
		return err
	}
	defer rows.Close()
	jsMaps := ""
	for rows.Next() {
		if err = rows.Scan(i.holder...); err != nil {
			return err
		}
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(i.name, i.holder, i.mode))
	}
	if len(jsMaps) < 3 {
		i.jsonList = "[]"
	} else {
		i.jsonList = fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	}
	return nil
}

func (i *Items) getRows() (*sql.Rows, error) {
	return db.Query(i.query)
}

func (i *Items) makeQuery() error {
	if err := testForExistingTable(i.name); err != nil {
		return err
	}
	var addNum int
	var sels, joins string
	if i.mode != "raw" {
		sels, joins, addNum = makeAddJoins(i.name)
	}

	i.query = fmt.Sprintf("SELECT %s.* %s FROM %s %s",
		i.name, sels, i.name, joins)

	if i.active == "unactive" {
		i.query = fmt.Sprintf("%s WHERE %s.is_active=0", i.query, i.name)
	} else if i.active != "all" {
		i.query = fmt.Sprintf("%s WHERE %s.is_active=1", i.query, i.name)
	}

	i.holder = makeHolder(len(Md.Models[i.name].Columns) + addNum)
	return nil
}

type FilteredItems struct {
	Items
	filterColumn string
	operator     string
	value        string
}

func (f *FilteredItems) Get() error {
	if err := f.makeQuery(); err != nil {
		return err
	}
	rows, err := f.getRows()
	if err != nil {
		return err
	}
	defer rows.Close()
	jsMaps := ""
	for rows.Next() {
		if err = rows.Scan(f.holder...); err != nil {
			return err
		}
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(f.name, f.holder, f.mode))
	}
	if len(jsMaps) < 3 {
		f.jsonList = "[]"
	} else {
		f.jsonList = fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	}

	return nil
}

func (f *FilteredItems) makeQuery() error {
	operator, ok := operators[f.operator]
	if !ok {
		return fmt.Errorf("operator '%s' does not exist", f.operator)
	}
	if err := testForExistingColumn(f.name, f.filterColumn); err != nil {
		return err
	}
	if err := f.Items.makeQuery(); err != nil {
		return err
	}
	field_type := Md.Models[f.name].Model[f.filterColumn].Type
	op_str := ""
	if field_type == "str" {
		op_str = fmt.Sprintf("LOW(%s.%s) %s ?", f.name, f.filterColumn, operator)
	} else {
		op_str = fmt.Sprintf("%s.%s %s ?", f.name, f.filterColumn, operator)
	}

	if f.active == "all" {
		f.query = fmt.Sprintf("%s WHERE %s", f.query, op_str)
	} else {
		f.query = fmt.Sprintf("%s AND %s", f.query, op_str)
	}
	fmt.Println(f.query)
	return nil
}

func (f *FilteredItems) getRows() (*sql.Rows, error) {
	return query(Md.Models[f.name].Model[f.filterColumn].Type, f.value, f.query)
}

type BetweenItems struct {
	Items
	filterColumn string
	value_start  string
	value_end    string
}

func (b *BetweenItems) Get() error {
	if err := b.makeQuery(); err != nil {
		return err
	}
	rows, err := b.getRows()
	if err != nil {
		return err
	}
	defer rows.Close()
	jsMaps := ""
	for rows.Next() {
		if err = rows.Scan(b.holder...); err != nil {
			return err
		}
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(b.name, b.holder, b.mode))
	}
	if len(jsMaps) < 3 {
		b.jsonList = "[]"
	} else {
		b.jsonList = fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	}

	return nil
}

func (b *BetweenItems) makeQuery() error {
	if err := testForExistingColumn(b.name, b.filterColumn); err != nil {
		return err
	}
	if err := b.Items.makeQuery(); err != nil {
		return err
	}
	if b.active == "all" {
		b.query = fmt.Sprintf("%s WHERE %s.%s BETWEEN ? and ?", b.query, b.name, b.filterColumn)
	} else {
		b.query = fmt.Sprintf("%s AND %s.%s BETWEEN ? and ?", b.query, b.name, b.filterColumn)
	}
	return nil
}

func (b *BetweenItems) getRows() (*sql.Rows, error) {
	return queryBetween(Md.Models[b.name].Model[b.filterColumn].Type, b.value_start, b.value_end, b.query)
}

func GetSum(table, column string) (float64, error) {
	if err := testForExistingTable(table); err != nil {
		return 0.0, err
	}
	if err := testForExistingColumn(table, column); err != nil {
		return 0.0, err
	}
	query := fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE is_active=1", column, table)
	row := db.QueryRow(query)
	var sum float64
	err := row.Scan(&sum)
	return sum, err
}

func GetSumFilter(table, sum_column, filter_column, operator, value string) (float64, error) {
	if err := testForExistingTable(table); err != nil {
		return 0.0, err
	}
	if err := testForExistingColumn(table, sum_column); err != nil {
		return 0.0, err
	}
	if err := testForExistingColumn(table, filter_column); err != nil {
		return 0.0, err
	}
	operator, ok := operators[operator]
	if !ok {
		return 0.0, fmt.Errorf("operator '%s' does not exist", operator)
	}
	query := fmt.Sprintf("SELECT SUM(%s) FROM %s WHERE %s %s %s AND is_active=1",
		sum_column, table, filter_column, operator, value)
	row := db.QueryRow(query)
	var sum float64
	err := row.Scan(&sum)
	return sum, err
}

func UpdateItem(table string, body io.ReadCloser) error {
	if err := testForExistingTable(table); err != nil {
		return err
	}

	req := map[string]interface{}{}
	decoder := json.NewDecoder(body)
	defer body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		return err
	}

	params_len := len(Md.Models[table].Columns)
	params := make([]interface{}, params_len)
	sql := ""
	for i, column := range Md.Models[table].Columns[1:] {
		params[i] = req[column]
		sql += fmt.Sprintf(", %s=?", column)
	}
	sql = fmt.Sprintf("UPDATE %s SET %s WHERE id=?", table, sql[2:])
	params[params_len-1] = req["id"]
	_, err = db.Exec(sql, params...)
	if err != nil {
		return err
	}
	return nil
}

func CreateItem(table string, body io.ReadCloser) error {
	if err := testForExistingTable(table); err != nil {
		return err
	}

	req := map[string]interface{}{}
	decoder := json.NewDecoder(body)
	defer body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		return err
	}

	params_len := len(Md.Models[table].Columns)
	params := make([]interface{}, params_len)
	sql := ""
	vals := ""
	for i, column := range Md.Models[table].Columns[1:] {
		params[i] = req[column]
		sql += fmt.Sprintf(", %s", column)
		vals += ", ?"
	}
	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, sql[2:], vals[2:])
	params[params_len-1] = req["id"]
	_, err = db.Exec(sql, params...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItem(table string, id int, mode string) error {
	if err := testForExistingTable(table); err != nil {
		return err
	}
	var sql string
	if mode == "delete" {
		sql = fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET is_active = 0 WHERE id = ?", table)
	}
	_, err := db.Exec(sql, id)
	if err != nil {
		return err
	}
	return nil
}

type Base struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateBase(body io.ReadCloser) error {
	base := Base{}
	decoder := json.NewDecoder(body)
	defer body.Close()
	err := decoder.Decode(&base)
	if err != nil {
		return err
	}
	// var new_db *sql.DB
	dbs[base.Name], err = sql.Open("sqlite3", base.Name)
	if err != nil {
		return err
	}
	fmt.Println("Created base ", base.Name, base.Description)
	return nil
}

func GetModels() (string, error) {
	data, err := os.ReadFile("models.json")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// func GetModels() (Model, error) {
// 	res := Model{}
// 	data, err := os.ReadFile("models.json")
// 	if err != nil {
// 		return res, err
// 	}
// 	decoder := json.NewDecoder(strings.NewReader(string(data)))
// 	err = decoder.Decode(&res)
// 	if err != nil {
// 		r.Respond(nil, err)
// 	}
// 	r.Respond(res, nil)
// }
