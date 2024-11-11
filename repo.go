package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

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
	i.holder = makeHolder(tablesColumnsNum[i.name] + addNum)
	return nil
}

type Items struct {
	name     string
	mode     string
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
	i.query = fmt.Sprintf("SELECT %s.* %s FROM %s %s WHERE %s.is_active=1",
		i.name, sels, i.name, joins, i.name)
	i.holder = makeHolder(tablesColumnsNum[i.name] + addNum)
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
	f.query = fmt.Sprintf("%s AND %s.%s %s ?", f.query, f.name, f.filterColumn, operator)
	return nil
}

func (f *FilteredItems) getRows() (*sql.Rows, error) {
	return query(tablesColumnsRawMap[f.name][f.filterColumn].ctype, f.value, f.query)
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

	params_len := len(tablesColumnsRaw[table])
	params := make([]interface{}, params_len)
	sql := ""
	for i, column := range tablesColumnsRaw[table][1:] {
		params[i] = req[column.cname]
		sql += fmt.Sprintf(", %s=?", column.cname)
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

	params_len := len(tablesColumnsRaw[table])
	params := make([]interface{}, params_len)
	sql := ""
	vals := ""
	for i, column := range tablesColumnsRaw[table][1:] {
		params[i] = req[column.cname]
		sql += fmt.Sprintf(", %s", column.cname)
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
