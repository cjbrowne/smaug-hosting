package helpers

import (
	"reflect"
	"strings"
)

func GetColsForType(i interface{}) []string {
	t := reflect.TypeOf(i)
	cols := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		col := t.Field(i).Tag.Get("db")
		if col == "" {
			col = strings.ToLower(t.Field(i).Name)
		}
		if col != "-" {
			cols = append(cols, col)
		}
	}
	return cols
}
