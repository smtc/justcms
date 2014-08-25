package utils

import "testing"

type TestJsonModel struct {
	Name     string
	Age      int
	Password string
}

func TestToJson(t *testing.T) {
	m := TestJsonModel{
		Name: "name",
		Age:  18,
	}

	parms := []string{"Name", "Age", "None"}

	s, err := ToJson(m, parms, FilterModeInclude)
	if err != nil {
		t.Fatal(err.Error())
	}

	if s != "{\"Age\":18,\"Name\":\"name\"}" {
		t.Error("error")
	}

	s, err = ToJson(m, parms, FilterModeIncludeMust)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"Age\":18,\"Name\":\"name\",\"None\":null}" {
		t.Error("error")
	}

	s, err = ToJson(m, []string{"Password"}, FilterModeExclude)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"Age\":18,\"Name\":\"name\"}" {
		t.Error("error")
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
