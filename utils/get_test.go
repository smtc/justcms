package utils

import "testing"

func TestGetValue(t *testing.T) {
	test := func(v interface{}) {
		age := GetValue(v, "Age").Interface().(int)
		if age != 1 {
			t.Fail()
		}

		age = GetInt(v, "Age", 0)
		if age != 1 {
			t.Fail()
		}

		age = GetInt(v, "age", 0)
		if age != 0 {
			t.Fail()
		}
	}

	type User struct {
		Age int
	}

	user := User{
		Age: 1,
	}

	test(user)

	m := map[string]interface{}{
		"Age": 1,
	}

	test(m)
}
