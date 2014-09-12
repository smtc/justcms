package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

var (
	// reserve table name
	NotAllowTables = []string{"user", "role", "post"}
)

type Table struct {
	Id        int64     `json:"id"`
	Name      string    `sql:"size:45;not null;unique" json:"name"`
	Alias     string    `sql:"size:45" json:"alias"`
	Des       string    `Sql:"size:512" json:"des"`
	CreatedAt time.Time `json:"created_at"`
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
	db.Model(Table{}).Where("name=?", t.Name).Count(&count)
	return count > 0
}

func (t *Table) Get(id int64) error {
	db := getTableDB()
	return db.First(t, id).Error
}

func (t *Table) Save() error {
	db := getTableDB()
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
