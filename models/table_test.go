package models

import "testing"

func TestTable(t *testing.T) {
	table := Table{
		Name:  "test",
		Alias: "测试",
		Des:   "测试表",
	}

	if err := table.Save(); err != nil {
		t.Fatal(err.Error())
	}

	var table2 Table
	if err := table2.Get(table.Id); err != nil {
		t.Fatal(err.Error())
	}

	table.Id = 0

	if err := table.Save(); err == nil {
		t.Fatal("table should be existed")
	}

	table.Name = "table"
	if err := table.Save(); err != nil {
		t.Fatal(err.Error())
	}

	if table.Name == table2.Name {
		t.Fatal(table.Name + "==" + table2.Name)
	}

	table3 := Table{
		Name:  "test",
		Alias: "用户",
		Columns: []Column{
			Column{Name: "id", Alias: "id", Type: "int", Size: 20},
			Column{Name: "name", Alias: "姓名", Type: "string", Size: 45},
			Column{Name: "age", Alias: "年龄", Type: "int", Size: 11},
		},
	}
	exist := table3.Exist()
	if !exist {
		t.Fatal("table test should be existed")
	}
	if err := table3.Save(); err == nil {
		t.Fatal("table3 save should error")
	}
	table3.Name = "user"
	if err := table3.Save(); err != nil {
		t.Fatal(err.Error())
	}

	table2.Delete()
	TableDelete("name='table'")

	// column test
	c0 := table3.Columns[0]
	if c0.Name != "id" {
		t.Fatal("column should be 'id', not ", c0.Name)
	}

	c1 := table3.Field("name")
	if c1.Alias != "姓名" {
		t.Fatal("column should be '姓名', not ", c1.Alias)
	}

	c2 := Column{TableId: table3.Id, Name: "email", Alias: "邮箱", Type: "string", Size: 100}
	c2.Save()
	if err := table3.GetColumns(); err != nil {
		t.Fatal(err.Error())
	}
	if len(table3.Columns) != 4 {
		t.Fatal("table3's Column length should be 4 ", len(table3.Columns))
	}

	if err := table3.Delete(); err != nil {
		t.Fatal(err.Error())
	}
}
