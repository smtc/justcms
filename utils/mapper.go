package utils

import (
	"log"
	"reflect"
	"strings"
	"time"
)

type filterMode int

const (
	FilterModeInclude filterMode = iota
	FilterModeIncludeMust
	FilterModeExclude
)

func filterMap(m map[string]interface{}, keys []string, mode filterMode) map[string]interface{} {
	var temp map[string]interface{}
	if mode == FilterModeExclude {
		temp = m
	} else {
		temp = make(map[string]interface{})
	}

	for _, k := range keys {
		k = strings.ToLower(k)
		if mode == FilterModeExclude {
			delete(temp, k)
		} else {
			val, ok := m[k]
			if ok || mode == FilterModeIncludeMust {
				temp[k] = val
			}
		}
	}
	return temp
}

func ToMap(v interface{}, keys []string, mode filterMode) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	fv := reflect.ValueOf(v)
	//ft := reflect.TypeOf(v)
	switch fv.Kind() {
	case reflect.Map:
		for _, k := range fv.MapKeys() {
			value := replaceTime(fv.MapIndex(k).Interface())
			m[k.String()] = value
		}

	case reflect.Struct:
		for i := 0; i < fv.NumField(); i++ {
			typeField := fv.Type().Field(i)
			value := replaceTime(fv.Field(i).Interface())
			tag := typeField.Tag.Get("json")
			if tag == "-" {
				continue
			}
			if tag == "" {
				tag = strings.ToLower(typeField.Name)
			}
			m[tag] = value
		}

	}

	return filterMap(m, keys, mode), nil
}

func replaceTime(v interface{}) interface{} {
	if t, ok := v.(time.Time); ok {
		return Time{t, ""}
	}
	return v
}

func ToMapList(v interface{}, keys []string, mode filterMode) ([]map[string]interface{}, error) {
	s := reflect.ValueOf(v)

	if s.Kind() != reflect.Slice {
		m, err := ToMap(v, keys, mode)
		return []map[string]interface{}{m}, err
	}

	var (
		ms    map[string]interface{}
		err   error
		count = s.Len()
		list  = make([]map[string]interface{}, count)
	)

	for i := 0; i < count; i++ {
		ms, err = ToMap(s.Index(i).Interface(), keys, mode)
		if err != nil {
			log.Println(err.Error())
			return list, err
		}
		list[i] = ms
	}

	return list, nil
}

func ToMapOnly(v interface{}) (map[string]interface{}, error) {
	return ToMap(v, []string{}, FilterModeExclude)
}
