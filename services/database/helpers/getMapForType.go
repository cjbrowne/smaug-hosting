package helpers

import "reflect"

// Deprecated: this was always a stupid idea
func GetMapForType(i interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	cols := GetColsForType(i)
	for _, c := range cols {
		v := reflect.ValueOf(i).FieldByName(c)
		switch v.Kind() {

		}
	}
	return res
}
