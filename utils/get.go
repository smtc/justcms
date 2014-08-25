package utils

import "reflect"

func GetInterface(v interface{}, key string) interface{} {
	fv := reflect.ValueOf(v)
	switch fv.Kind() {
	case reflect.Struct:
		return nil
		break
	}

	return fv
}

func GetInt(v interface{}, key string, def int) int {
	return def
}
