package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

var (
	// reserve table name
	NotAllowTables = []string{"account", "role", "post"}
)

type Table struct {
	Id        int64     `json:"id"`
	Name      string    `sql:"size:45;not null;unique" json:"name"`
	Alias     string    `sql:"size:45" json:"alias"`
	Des       string    `Sql:"size:512" json:"des"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
	Columns   []Column  `json:"columns"`
}

func getTableDB() *gorm.DB {
	return database.GetDB(dynamic_db)
}

func (t *Table) Exist() bool {
	for _, nat := range NotAllowTables {
		if t.Name == nat {
			return true
		}
	}
	db := getTableDB()
	var count int
	db.Model(Table{}).Where("id!=? and name=?", t.Id, t.Name).Count(&count)
	return count > 0
}

func (t *Table) Get(id int64) error {
	db := getTableDB()
	return db.First(t, id).Error
}

func (t *Table) GetColumns() error {
	db := getTableDB()
	t.Columns = nil
	return db.Model(t).Related(&t.Columns).Error
}

func (t *Table) Refresh() error {
	db := getTableDB()
	err := db.First(t, t.Id).Error
	if err != nil {
		return err
	}
	return t.GetColumns()
}

func (t *Table) Save() error {
	db := getTableDB()
	if t.Exist() {
		return fmt.Errorf("table '%v' is existed", t.Name)
	}
	if t.Id == 0 {
		t.CreatedAt = time.Now()
	}
	t.EditAt = time.Now()
	return db.Save(t).Error
}

func (t *Table) Delete() error {
	db := getTableDB()
	return db.Delete(t).Error
}

func TableDelete(where string) {
	db := getTableDB()
	db.Where(where).Delete(&Table{})
}

func TableList() ([]Table, error) {
	db := getTableDB()
	var tbls []Table

	err := db.Find(&tbls).Error
	return tbls, err
}

func (t *Table) Field(name string) Column {
	var column Column
	for _, c := range t.Columns {
		if c.Name == name {
			column = c
		}
	}
	return column
}
