package utils

import (
	"encoding/json"
	"log"
	"reflect"
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
	obj, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(obj, &m)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return filterMap(m, keys, mode), nil
}

func ToMapList(v interface{}, keys []string, mode filterMode) ([]map[string]interface{}, error) {
	s := reflect.ValueOf(v)
	var (
		ms    map[string]interface{}
		err   error
		count = s.Len()
	)
	list := make([]map[string]interface{}, count)
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
