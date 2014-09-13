/*
* get value from struct or map by string key
* if key does not exist, return given default value
 */

package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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

func GetInterface(v interface{}, key string, def interface{}) interface{} {
	val := GetValue(v, key)
	if !val.IsValid() {
		return def
	}

	return val.Interface()
}

func GetInt64(v interface{}, key string, def int64) int64 {
	itf := GetInterface(v, key, def)
	if i, ok := itf.(int64); ok {
		return i
	}
	if i, ok := itf.(int); ok {
		return int64(i)
	}
	if ss, ok := itf.([]string); ok {
		itf = ss[0]
	}
	if s, ok := itf.(string); ok {
		if i, err := strconv.ParseInt(s, 0, 64); err == nil {
			return i
		}
	}
	return def
}

func GetInt(v interface{}, key string, def int) int {
	i := GetInt64(v, key, int64(def))
	return int(i)
}

func GetFloat64(v interface{}, key string, def float64) float64 {
	itf := GetInterface(v, key, def)
	if f, ok := itf.(float64); ok {
		return f
	}
	if ss, ok := itf.([]string); ok {
		itf = ss[0]
	}
	if s, ok := itf.(string); ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return def
}

func GetFloat32(v interface{}, key string, def float32) float32 {
	f := GetFloat64(v, key, float64(def))
	return float32(f)
}

func GetString(v interface{}, key string, def string) string {
	itf := GetInterface(v, key, def)
	if s, ok := itf.(string); ok {
		return s
	}
	if ss, ok := itf.([]string); ok {
		return strings.Join(ss, ",")
	}
	if t, ok := itf.(time.Time); ok {
		return t.Format(TIMEFORMAT)
	}

	return fmt.Sprintf("%v", itf)
}

func GetTime(v interface{}, key string, def time.Time, fmt string) time.Time {
	itf := GetInterface(v, key, def)
	if t, ok := itf.(time.Time); ok {
		return t
	}
	if ss, ok := itf.([]string); ok {
		itf = ss[0]
	}
	if s, ok := itf.(string); ok {
		if t, err := time.Parse(fmt, s); err == nil {
			return t
		}
	}

	return def
}

func GetBool(v interface{}, key string, def bool) bool {
	itf := GetInterface(v, key, def)
	if b, ok := itf.(bool); ok {
		return b
	}
	if i, ok := itf.(int); ok {
		return i > 0
	}
	if ss, ok := itf.([]string); ok {
		itf = ss[0]
	}
	if s, ok := itf.(string); ok {
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
	}

	return def
}

func GetBytes(v interface{}, key string, def []byte) []byte {
	itf := GetInterface(v, key, def)
	if b, ok := itf.([]byte); ok {
		return b
	}
	if ss, ok := itf.([]string); ok {
		itf = ss[0]
	}
	if s, ok := itf.(string); ok {
		return []byte(s)
	}

	return def
}
