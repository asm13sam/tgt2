package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	INT = iota
	STRING
	BOOL
	FLOAT
)

var db *sql.DB

func DBconnect(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	return err
}

func DbClose() error {
	return db.Close()
}

func testForExistingTable(tableName string) bool {
	_, ok := tablesColumns1[tableName]
	return ok
}

func testForExistingColumn(tableName, columnName string) bool {
	_, ok := tablesColumns1[tableName][columnName]
	return ok
}

func makeSql(tableName, operator string, params ...string) (string, error) {
	var res string
	var err error
	if !testForExistingTable(tableName) {
		return res, fmt.Errorf("table %s does not exist", tableName)
	}
	switch operator {
	case "by_id":
		res = fmt.Sprintf("SELECT * FROM %s WHERE id=?", tableName)
	case "all":
		res = fmt.Sprintf("SELECT * FROM %s WHERE is_active=1", tableName)
	case "filter":
		op, ok := operands[params[1]]
		if !ok {
			return res, fmt.Errorf("operand '%s' does not exist", params[1])
		}
		if !testForExistingColumn(tableName, params[0]) {
			return res, fmt.Errorf("column '%s' in table '%s' does not exist", params[0], tableName)
		}
		res = fmt.Sprintf("SELECT * FROM %s WHERE is_active=1 AND %s %s ?", tableName, params[0], op)
	default:
		err = fmt.Errorf("operator %s does not exist", operator)
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
