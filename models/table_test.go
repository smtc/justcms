package models

import (
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	table := Table{
		Name:      "test",
		Alias:     "测试",
		Des:       "测试表",
		CreatedAt: time.Now(),
	}

	if err := table.Save(); err != nil {
		t.Fatal(err.Error())
	}

	var table2 Table
	if err := table2.Get(table.Id); err != nil {
		t.Fatal(err.Error())
	}

	table.Id = 0
	table.Name = "table"
	if err := table.Save(); err != nil {
		t.Fatal(err.Error())
	}

	if table.Name == table2.Name {
		t.Fatal(table.Name + "==" + table2.Name)
	}

	table3 := Table{Name: "test"}
	exist := table3.Exist()
	if !exist {
		t.Fatal("table test should exist")
	}
	if err := table3.Save(); err == nil {
		t.Fatal("table3 save should error")
	}

	//table2.Delete()
	//TableDelete("name='名字'")

}
