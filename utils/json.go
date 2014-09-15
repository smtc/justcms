package utils

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
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

func ToJsonOnly(v interface{}) (string, error) {
	m, err := ToMapOnly(v)
	if err != nil {
		return "", err
	}

	obj, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(obj), nil
}

func ToStruct(data []byte, v interface{}) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	var getTime = func(itf interface{}) time.Time {
		if t, ok := itf.(time.Time); ok {
			return t
		}
		if s, ok := itf.(string); ok {
			if t, err := time.Parse(TIMEFORMAT, s); err == nil {
				return t
			}
		}
		return TIMEDEFAULT
	}

	fv := reflect.ValueOf(v).Elem()
	var field reflect.Value
	for i := 0; i < fv.NumField(); i++ {
		field = fv.Field(i)
		if field.Kind() == reflect.Struct && field.Type() == timeType {
			typeField := fv.Type().Field(i)
			tag := typeField.Tag.Get("json")
			if tag == "" {
				tag = strings.ToLower(typeField.Name)
			}
			m[tag] = getTime(m[tag])
		}
	}

	if bs, err := json.Marshal(m); err == nil {
		err = json.Unmarshal(bs, v)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
