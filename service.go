package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	sqlite "github.com/mattn/go-sqlite3"
)

const (
	INT = iota
	STRING
	BOOL
	FLOAT
)

var db *sql.DB

var dbs = map[string]*sql.DB{}

func low(s string) string {
	return strings.ToLower(s)
}

func DBconnect(dbFile string) error {
	var err error
	sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("LOW", low, true); err != nil {
				return err
			}
			return nil
		},
	})

	db, err = sql.Open("sqlite3_custom", dbFile)
	return err
}

func DbClose() error {
	return db.Close()
}

func testForExistingTable(tableName string) error {
	if _, ok := Md.Models[tableName]; !ok {
		return fmt.Errorf("table %s does not exist", tableName)
	}
	return nil
}

func testForExistingColumn(tableName, columnName string) error {
	if _, ok := Md.Models[tableName].Model[columnName]; !ok {
		return fmt.Errorf("column '%s' in table '%s' does not exist", columnName, tableName)
	}
	return nil
}

func makeAddJoins(table string) (addSelect, addJoins string, addNum int) {
	for _, col := range Md.Models[table].Columns {
		joinTable, ok := strings.CutSuffix(col, "_id")
		if ok {
			addNum++
			origTable, ok := strings.CutSuffix(joinTable, "2")
			if ok {
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, joinTable)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s AS %s ON %s.%s = %s.id`,
					origTable, joinTable, table, col, joinTable,
				)
			} else if joinTable == table {
				jtShort := joinTable[:2]
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, jtShort)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s AS %s ON %s.%s = %s.id`,
					joinTable, jtShort, table, col, jtShort,
				)
			} else {
				addSelect += fmt.Sprintf(`, IFNULL(%s.name, "")`, joinTable)
				addJoins += fmt.Sprintf(
					` LEFT JOIN %s ON %s.%s = %s.id`,
					joinTable, table, col, joinTable,
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
	var t ItemColumn
	columns := Md.Models[tableName].Columns
	if mode != "raw" {
		columns = append(columns, Md.Models[tableName].WColumns...)
	}

	for i, s := range holder {
		v := *s.(*string)
		var ok bool
		if t, ok = Md.Models[tableName].Model[columns[i]]; !ok {
			t = Md.Models[tableName].WModel[columns[i]]
		}
		if t.Type == "bool" {
			if v == "0" {
				v = "false"
			} else {
				v = "true"
			}
		} else if t.Type == "str" {
			if strings.Contains(v, "\\") {
				fmt.Println(v)
				v = strings.ReplaceAll(v, "\\", "/")
			}
			if strings.Contains(v, "\"") {
				fmt.Println(v)
				v = strings.ReplaceAll(v, "\"", "'")
			}

		}
		js += fmt.Sprintf(typesTemplates[t.Type], columns[i], v)
	}
	return fmt.Sprintf("{\n%s\n}", js[:len(js)-2])
}

func query(ctype, value, sql string) (*sql.Rows, error) {
	switch ctype {
	case "int":
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return db.Query(sql, num)
	case "str":
		return db.Query(sql, value)
	case "bool":
		return db.Query(sql, value == "true" || value == "1")
	case "float":
		fnum, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return db.Query(sql, fnum)
	default:
		return nil, fmt.Errorf("unknown type code %s for value %s", ctype, value)
	}
}

func queryBetween(ctype, value_start, value_end, sql string) (*sql.Rows, error) {
	switch ctype {
	case "int":
		num1, err := strconv.Atoi(value_start)
		if err != nil {
			return nil, err
		}
		num2, err := strconv.Atoi(value_end)
		if err != nil {
			return nil, err
		}
		return db.Query(sql, num1, num2)
	case "str":
		return db.Query(sql, value_start, value_end)
	case "bool":
		return db.Query(sql, value_start == "true" || value_start == "1", value_end == "true" || value_end == "1")
	case "float":
		fnum1, err := strconv.ParseFloat(value_start, 64)
		if err != nil {
			return nil, err
		}
		fnum2, err := strconv.ParseFloat(value_end, 64)
		if err != nil {
			return nil, err
		}
		return db.Query(sql, fnum1, fnum2)
	default:
		return nil, fmt.Errorf("unknown type code %s for values %s, %s", ctype, value_start, value_end)
	}
}
