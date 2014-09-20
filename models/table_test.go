package models

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func TestTable(t *testing.T) {
	var err error
	table := Table{
		Name:  "test",
		Alias: "测试",
		Des:   "测试表",
	}

	if err = table.Save(); err != nil {
		t.Fatal(err.Error())
	}

	table2, err := GetTable(table.Id)
	if err != nil {
		t.Fatal(err.Error())
	}

	table.Id = 0

	if err = table.Save(); err == nil {
		t.Fatal("table should be existed")
	}

	/*
		table.Name = "table"
		if err = table.Save(); err != nil {
			t.Fatal(err.Error())
		}

		if table.Name == table2.Name {
			t.Fatal(table.Name + "==" + table2.Name)
		}
	*/

	table3 := Table{
		Name:  "test",
		Alias: "用户",
		Columns: []Column{
			Column{Name: "id", Alias: "id", Type: AUTO_INCREMENT, Size: 20},
			Column{Name: "name", Alias: "姓名", Type: VARCHAR, Size: 45},
			Column{Name: "age", Alias: "年龄", Type: INT, Size: 11},
			Column{Name: "birthday", Alias: "生日", Type: DATE},
		},
	}
	exist := table3.Exist()
	if !exist {
		t.Fatal("table test should be existed")
	}
	if err = table3.Save(); err == nil {
		t.Fatal("table3 save should error")
	}
	table3.Name = "user"
	if err = table3.Save(); err != nil {
		t.Fatal(err.Error())
	}

	// column test
	c0 := table3.Columns[0]
	if c0.Name != "id" {
		t.Fatal("column should be 'id', not ", c0.Name)
	}

	c1 := table3.Field("name")
	if c1.Alias != "姓名" {
		t.Fatal("column should be '姓名', not ", c1.Alias)
	}

	c2 := Column{TableId: table3.Id, Name: "email1", Alias: "邮箱", Type: VARCHAR, Size: 100}
	c2.Save()
	if err = table3.Refresh(); err != nil {
		t.Fatal(err.Error())
	}

	c2.Name = "email"
	c2.Type = BIGINT
	c2.Save()

	if len(table3.Columns) != 5 {
		t.Fatal("table3's Column length should be 5 ", len(table3.Columns))
	}

	if err = c2.Delete(); err != nil {
		t.Fatal(err.Error())
	}

	table3.Refresh()
	if len(table3.Columns) != 4 {
		t.Fatal("table3's Column length should be 4 ", len(table3.Columns))
	}

	if err = ColumnDelete("name like ?", "%e%"); err != nil {
		t.Fatal(err.Error())
	}
	table3.Refresh()
	if len(table3.Columns) != 2 {
		t.Fatal("table3's Column length should be 2 ", len(table3.Columns))
	}

	db := GetDB(DYNAMIC_DB)
	scope := db.NewScope(nil)
	suc := gorm.NewDialect("mysql").HasTable(scope, "user")
	if !suc {
		t.Fatal("table user not exist")
	}

	stt := Struct{}
	stt.GetStruct(&table3)
	//json, err := goutils.ToJsonOnly(stt)
	//println(json)

	table2.Delete()
	table.Delete()
	if err = table3.Delete(); err != nil {
		t.Fatal(err.Error())
	}
}
