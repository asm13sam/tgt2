package main

import (
	"database/sql"
	"fmt"

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
