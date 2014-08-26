package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGetValue(t *testing.T) {
	var (
		birthday time.Time
		now      time.Time
		keys     []byte
	)
	birthday, _ = time.Parse(TIMEFORMAT, "2012-03-04 05:06:07")
	keys = []byte("123456")
	now = time.Now()

	test := func(v interface{}) {
		var (
			intVal  int
			strVal  string
			f64Val  float64
			f32Val  float32
			timeVal time.Time
			boolVal bool
			byteVal []byte
		)

		if intVal = GetValue(v, "Age").Interface().(int); intVal != 1 {
			t.Fatal(fmt.Sprintf("Age should be 1, not %v", intVal))
		}

		if intVal = GetInt(v, "Age", 0); intVal != 1 {
			t.Fatal(fmt.Sprintf("Age should be 1, not %v", intVal))
		}

		if intVal = GetInt(v, "age", 0); intVal != 0 {
			t.Fatal(fmt.Sprintf("age should be 0, not %v", intVal))
		}

		if intVal = GetInt(v, "StrAge", 0); intVal != 100 {
			t.Fatal(fmt.Sprintf("StrAge should be 100, not %v", intVal))
		}

		if strVal = GetValue(v, "StrAge").Interface().(string); strVal != "100" {
			t.Fatal(fmt.Sprintf("StrAge should be '100', not %v", strVal))
		}

		if strVal = GetString(v, "StrAge", ""); strVal != "100" {
			t.Fatal(fmt.Sprintf("StrAge should be '100', not %v", strVal))
		}

		if strVal = GetString(v, "Empty", ""); strVal != "" {
			t.Fatal(fmt.Sprintf("Empty should be empty, not %v", strVal))
		}

		if strVal = GetString(v, "Age", ""); strVal != "1" {
			t.Fatal(fmt.Sprintf("Age should be '1', not %v", strVal))
		}

		if strVal = GetString(v, "Birthday", ""); strVal != "2012-03-04 05:06:07" {
			t.Fatal(fmt.Sprintf("Age should be '2012-03-04 05:06:07', not %v", strVal))
		}

		if f64Val = GetFloat64(v, "Money", 0); f64Val != 123.45 {
			t.Fatal(fmt.Sprintf("StrMoney should be 123.45, not %v", f64Val))
		}

		if f64Val = GetFloat64(v, "StrMoney", 0); f64Val != 678.90 {
			t.Fatal(fmt.Sprintf("StrMoney should be 678.90, not %v", f64Val))
		}

		if f64Val = GetFloat64(v, "money", 0); f64Val != 0 {
			t.Fatal(fmt.Sprintf("money should be 0, not %v", f64Val))
		}

		if f32Val = GetFloat32(v, "StrMoney", 0); f32Val != 678.90 {
			t.Fatal(fmt.Sprintf("StrMoney should be 678.90, not %v", f64Val))
		}

		if timeVal = GetTime(v, "Birthday", now, TIMEFORMAT); timeVal != birthday {
			t.Fatal(fmt.Sprintf("Birthday should be 2012-03-04 05:06:07, not %v", timeVal))
		}

		if timeVal = GetTime(v, "StrBirthday", now, TIMEFORMAT); timeVal != birthday {
			t.Fatal(fmt.Sprintf("StrBirthday should be 2012-03-04 05:06:07, not %v", timeVal))
		}

		if timeVal = GetTime(v, "NoneBirthday", now, TIMEFORMAT); timeVal != now {
			t.Fatal(fmt.Sprintf("NoneBirthday should be %v, not %v", now, timeVal))
		}

		if boolVal = GetBool(v, "Active", false); !boolVal {
			t.Fatal(fmt.Sprintf("Active should be true, not %v", boolVal))
		}

		if boolVal = GetBool(v, "StrActive", false); !boolVal {
			t.Fatal(fmt.Sprintf("StrActive should be true, not %v", boolVal))
		}

		if boolVal = GetBool(v, "Age", false); !boolVal {
			t.Fatal(fmt.Sprintf("Age should be true, not %v", boolVal))
		}

		if boolVal = GetBool(v, "None", false); boolVal {
			t.Fatal(fmt.Sprintf("None should be false, not %v", boolVal))
		}

		// 暂时先不判断slice相等, 只判断长度
		if byteVal = GetBytes(v, "Keys", []byte{}); len(byteVal) != len(keys) {
			t.Fatal(fmt.Sprintf("Keys should be [49 50 51 52 53 54], not %v", byteVal))
		}

		if byteVal = GetBytes(v, "StrKeys", []byte{}); len(byteVal) != len(keys) {
			t.Fatal(fmt.Sprintf("StrKeys should be [49 50 51 52 53 54], not %v", byteVal))
		}

		if byteVal = GetBytes(v, "None", []byte{}); len(byteVal) != 0 {
			t.Fatal(fmt.Sprintf("None should be [], not %v", byteVal))
		}
	}

	// ======================================================

	type User struct {
		Age         int
		StrAge      string
		Birthday    time.Time
		StrBirthday string
		Money       float64
		StrMoney    string
		Active      bool
		StrActive   string
		Keys        []byte
		StrKeys     string
	}

	user := User{
		Age:         1,
		StrAge:      "100",
		Birthday:    birthday,
		StrBirthday: "2012-03-04 05:06:07",
		Money:       123.45,
		StrMoney:    "678.90",
		Active:      true,
		StrActive:   "true",
		Keys:        keys,
		StrKeys:     "123456",
	}

	test(user)

	// ======================================================

	m := map[string]interface{}{
		"Age":         1,
		"StrAge":      "100",
		"Birthday":    birthday,
		"StrBirthday": "2012-03-04 05:06:07",
		"Money":       123.45,
		"StrMoney":    "678.90",
		"Active":      true,
		"StrActive":   "TRUE",
		"Keys":        keys,
		"StrKeys":     "123456",
	}

	test(m)
}
