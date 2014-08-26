package utils

import (
	"testing"
	"time"
)

func TestToMap(t *testing.T) {
	type TestMapModel struct {
		Name     string
		Age      int
		Password string
		Birthday Time
	}

	list := []TestMapModel{
		TestMapModel{Name: "test1", Age: 18},
		TestMapModel{Name: "test2", Age: 19, Birthday: Time{time.Now(), ""}},
	}

	if m, err := ToMapList(list, []string{}, FilterModeExclude); err != nil {
		t.Fatal(err.Error())
	} else {
		_, err := ToJson(m, []string{}, FilterModeExclude)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}
