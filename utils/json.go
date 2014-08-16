package utils

import (
	"encoding/json"
	"log"
)

type jsonMode int

const (
	jsonInclude jsonMode = iota
	jsonIncludeMust
	jsonExclude
)

func ToMap(v interface{}) (map[string]interface{}, error) {
	obj, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	mobj := make(map[string]interface{})
	err = json.Unmarshal(obj, &mobj)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return mobj, nil
}

func toJson(v interface{}, keys []string, mode jsonMode) (string, error) {
	mobj, err := ToMap(v)
	if err != nil {
		return "", err
	}

	var temp map[string]interface{}
	if mode == jsonExclude {
		temp = mobj
	} else {
		temp = make(map[string]interface{})
	}

	for _, k := range keys {
		if mode == jsonExclude {
			delete(temp, k)
		} else {
			val, ok := mobj[k]
			if ok || mode == jsonIncludeMust {
				temp[k] = val
			}
		}
	}

	obj, err := json.Marshal(temp)
	if err != nil {
		return "", err
	}

	return string(obj), nil

}

func ToJson(v interface{}, keys []string) (string, error) {
	return toJson(v, keys, jsonInclude)
}

func ToJsonMust(v interface{}, keys []string) (string, error) {
	return toJson(v, keys, jsonIncludeMust)
}

func ToJsonEx(v interface{}, keys []string) (string, error) {
	return toJson(v, keys, jsonExclude)
}
