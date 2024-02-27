package tool

import (
	"fmt"
	"reflect"
)

func CsvHeaders(s any) []string {
	t := reflect.TypeOf(s)
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		headers[i] = field.Tag.Get("csv")
	}
	return headers
}

func ToCsv(s any) []string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	row := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		row[i] = fmt.Sprintf("%v", field.Interface())
	}
	return row
}
