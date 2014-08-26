package utils

import "testing"

type TestMapModel struct {
	Name     string
	Age      int
	Password string
}

func TestToMap(t *testing.T) {
	list := []TestMapModel{
		TestMapModel{Name: "test1", Age: 18},
		TestMapModel{Name: "test2", Age: 19},
	}

	//for i := 0; i < 10000; i++ {
	ToMapList(list, []string{"Name"}, FilterModeInclude)
	//}
}
