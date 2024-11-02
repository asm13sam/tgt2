package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetItemRaw(tableName string, id int) (string, error) {
	sql, err := makeGetItemSqlRaw(tableName)
	if err != nil {
		return "", err
	}
	holder := makeHolder(tablesColumnsNum[tableName])
	row := db.QueryRow(sql, id)
	err = row.Scan(holder...)
	if err != nil {
		return "", err
	}
	jsMap := makeJsonMap(tableName, holder, "raw")
	return jsMap, nil
}

func GetItem(tableName string, id int) (string, error) {
	sql, addNum, err := makeGetItemSql(tableName)
	if err != nil {
		return "", err
	}
	fmt.Println(">>>>>>>", sql)
	holder := makeHolder(tablesColumnsNum[tableName] + addNum)
	row := db.QueryRow(sql, id)
	err = row.Scan(holder...)
	if err != nil {
		return "", err
	}
	jsMap := makeJsonMap(tableName, holder, "")
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
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(tableName, holder, "raw"))
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
	rows, err := query(tablesColumnsRawMap[tableName][filterColumn].ctype, value, sel)
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
		jsMaps += fmt.Sprintf("%s\n,", makeJsonMap(tableName, holder, "raw"))
	}
	if len(jsMaps) < 3 {
		return "[]", nil
	}
	jsList := fmt.Sprintf("[\n%s\n]", jsMaps[:len(jsMaps)-2])
	return jsList, nil
}
