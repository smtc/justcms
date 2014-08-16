package utils

import "testing"

type TestModel struct {
	Name     string
	Age      int
	Password string
}

func TestToJson(t *testing.T) {
	m := TestModel{
		Name: "name",
		Age:  18,
	}

	parms := []string{"Name", "Age", "None"}

	s, err := ToJson(m, parms)
	if err != nil {
		t.Fatal(err.Error())
	}

	if s != "{\"Age\":18,\"Name\":\"name\"}" {
		t.Error("error")
	}

	s, err = ToJsonMust(m, parms)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"Age\":18,\"Name\":\"name\",\"None\":null}" {
		t.Error("error")
	}

	s, err = ToJsonEx(m, []string{"Password"})
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "{\"Age\":18,\"Name\":\"name\"}" {
		t.Error("error")
	}

	for i := 0; i <= 10000; i++ {
		ToJson(m, parms)
	}

}
