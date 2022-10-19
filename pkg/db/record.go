package db

import (
	"fmt"
	"reflect"

	upperdb "github.com/upper/db/v4"
)

type Record struct {
	table   string
	columns []string
	values  [][]interface{}
}

func NewRecord(table string) *Record {
	return &Record{table: table}
}

func NewRecordWithObject(table string, object interface{}) *Record {
	type0 := reflect.TypeOf(object)
	switch type0.Kind() {
	case reflect.Ptr:
		return newRecordWithObject(table, []interface{}{object})
	case reflect.Slice:
		val := reflect.ValueOf(object)
		var objects []interface{}
		for i := 0; i < val.Len(); i++ {
			objects = append(objects, val.Index(i).Interface())
		}
		return newRecordWithObject(table, objects)
	default:
		panic(fmt.Sprintf("only support ptr and slice, not support:%s", type0.Kind()))
	}
}

func (r *Record) Value(column string, value interface{}) *Record {
	if len(r.values) == 0 {
		r.values = make([][]interface{}, 1)
	}
	r.columns = append(r.columns, column)
	r.values[0] = append(r.values[0], value)
	return r
}

func newRecordWithObject(table string, objects []interface{}) *Record {
	record := NewRecord(table)
	if len(objects) == 0 {
		return record
	}

	type0 := reflect.TypeOf(objects[0])
	if type0.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("object in record can't be %s", type0.Kind()))
	}
	ignoreColumnIdxs := map[int]interface{}{}
	for i := 0; i < type0.Elem().NumField(); i++ {
		dbTag := type0.Elem().Field(i).Tag.Get("db")
		if dbTag == "" {
			ignoreColumnIdxs[i] = struct{}{}
			continue
		}
		record.columns = append(record.columns, type0.Elem().Field(i).Tag.Get("db"))
	}

	record.values = make([][]interface{}, len(objects))
	for i := range objects {
		if reflect.TypeOf(objects[i]) != type0 {
			panic(fmt.Sprintf("object in record should be the same types %s", type0.String()))
		}
		value := reflect.ValueOf(objects[i]).Elem()
		for j := 0; j < value.NumField(); j++ {
			_, ignore := ignoreColumnIdxs[j]
			if ignore {
				continue
			}
			record.values[i] = append(record.values[i], value.Field(j).Interface())
		}
	}
	return record
}

func insertRecord(sqlBuilder upperdb.SQL, record *Record) error {
	smt := sqlBuilder.InsertInto(record.table).Columns(record.columns...)
	for i := range record.values {
		smt = smt.Values(record.values[i]...)
	}
	_, err := smt.Exec()
	return err
}
