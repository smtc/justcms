package utils

import (
	"testing"
	"time"
)

type TestJsonModel struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string
	Ignore   string `json:"-"`
	Birthday time.Time
}

func TestToJson(t *testing.T) {
	sbd := "2010-01-01 00:00:00"
	birthday, _ := time.Parse(TIMEFORMAT, sbd)
	m := TestJsonModel{
		Name:     "name",
		Age:      18,
		Birthday: birthday,
	}

	parms := []string{"name", "Age", "None"}
	var (
		s   string
		err error
	)

	s, err = ToJson(m, parms, FilterModeInclude)
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
	if s != "{\"age\":18,\"birthday\":\"2010-01-01 00:00:00\",\"name\":\"name\"}" {
		t.Error(s)
	}

	var newM TestJsonModel
	err = ToStruct([]byte(s), &newM)
	if err != nil {
		t.Fatal(err.Error())
	}
	if newM.Birthday != birthday {
		t.Fatal(newM.Birthday, "!=", birthday)
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
