package utils

import "reflect"

func GetValue(v interface{}, key string) reflect.Value {
	fv := reflect.ValueOf(v)
	switch fv.Kind() {

	case reflect.Struct:
		return fv.FieldByName(key)

	case reflect.Map:
		return fv.MapIndex(reflect.ValueOf(key))

	}

	return fv
}

func GetInt(v interface{}, key string, def int) int {
	if !GetValue(v, key).IsValid() {
		return def
	}

	i := GetValue(v, key).Interface().(int)
	/*
		if i, ok := GetValue(v, key).Interface().(int); ok {
			return i
		}
	*/
	return i
}
