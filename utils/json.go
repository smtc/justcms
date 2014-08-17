package utils

import (
	"encoding/json"
	"reflect"
)

func ToJson(v interface{}, keys []string, mode filterMode) (string, error) {
	s := reflect.ValueOf(v)
	var (
		m   interface{}
		err error
	)
	if s.Kind() == reflect.Slice {
		m, err = ToMapList(v, keys, mode)
	} else {
		m, err = ToMap(v, keys, mode)
	}
	if err != nil {
		return "", err
	}

	obj, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(obj), nil

}
