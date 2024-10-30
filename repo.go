package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetItem(tableName string, id int) (string, error) {
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

func GetItems(tableName string) (string, error) {
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
	if len(jsMaps) < 3 {
		return "[]", nil
	}
	jsList := fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	return jsList, nil
}
