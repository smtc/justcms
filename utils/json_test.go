package utils

import "testing"

type TestJsonModel struct {
	Name     string `json:"name"`
	Age      int
	Password string
	Ignore   string `json:"-"`
}

func TestToJson(t *testing.T) {
	m := TestJsonModel{
		Name: "name",
		Age:  18,
	}

	parms := []string{"name", "Age", "None"}

	s, err := ToJson(m, parms, FilterModeInclude)
	if err != nil {
		t.Fatal(err.Error())
	}

	if s != "{\"age\":18,\"name\":\"name\"}" {
		t.Error(s)
	}

	s, err = ToJson(m, parms, FilterModeIncludeMust)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"age\":18,\"name\":\"name\",\"none\":null}" {
		t.Error(s)
	}

	s, err = ToJson(m, []string{"Password"}, FilterModeExclude)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"age\":18,\"name\":\"name\"}" {
		t.Error(s)
	}

	list := []TestJsonModel{
		TestJsonModel{Name: "test1", Age: 18},
		TestJsonModel{Name: "test2", Age: 19},
	}

	s, err = ToJson(list, []string{}, FilterModeExclude)
	if err != nil {
		t.Fatal(err.Error())
	}
	//fmt.Println(s)
}
