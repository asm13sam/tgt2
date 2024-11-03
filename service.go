package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func testForExistingTable(tableName string) error {
	if _, ok := tablesColumnsRawMap[tableName]; !ok {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	return nil
}

func testForExistingColumn(tableName, columnName string) error {
	if _, ok := tablesColumnsRawMap[tableName][columnName]; !ok {
		return fmt.Errorf("column '%s' in table '%s' does not exist", columnName, tableName)
	}
	return nil
}

func makeAddJoins(table string) (addSelect, addJoins string, addNum int) {
	for _, col := range tablesColumnsRaw[table] {
		joinTable, ok := strings.CutSuffix(col.cname, "_id")
		if ok {
			addNum++
			origTable, ok := strings.CutSuffix(joinTable, "2")
			if ok {
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, joinTable)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s AS %s ON %s.%s = %s.id`,
					origTable, joinTable, table, col.cname, joinTable,
				)
			} else if joinTable == table {
				jtShort := joinTable[:2]
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, jtShort)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s AS %s ON %s.%s = %s.id`,
					joinTable, jtShort, table, col.cname, jtShort,
				)
			} else {
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, joinTable)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s ON %s.%s = %s.id`,
					joinTable, table, col.cname, joinTable,
				)
			}
		}
	}
	return
}

func makeHolder(number int) []interface{} {
	m := []interface{}{}
	for i := 0; i < number; i++ {
		s := ""
		m = append(m, &s)
	}
	return m
}

func makeJsonMap(tableName string, holder []interface{}, mode string) string {
	var js string
	var t int
	columns := tablesColumns[tableName]
	if mode == "raw" {
		columns = tablesColumnsRaw[tableName]
	}

	for i, s := range holder {
		v := *s.(*string)
		t = columns[i].ctype

		if t == BOOL {
			if v == "0" {
				v = "false"
			} else {
				v = "true"
			}
		}
		js += fmt.Sprintf(typesTemplates[t], columns[i].cname, v)
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
