package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
)

var (
	// reserve table name
	NotAllowTables = []string{}
)

type Table struct {
	Id        int64
	Name      string    `sql:"size:45;not null;unique"`
	Alias     string    `sql:"size:45"`
	Des       string    `sql:"size:512"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
	Columns   []Column
}

func getTableDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
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
	return db.Model(t).Related(&t.Columns).Order("name").Error
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
	var (
		db     = getTableDB()
		isNew  = false
		column Column
	)
	if t.Exist() {
		return fmt.Errorf("table '%v' is existed", t.Name)
	}
	if t.Id == 0 {
		t.CreatedAt = time.Now()
		isNew = true
	}
	t.EditAt = time.Now()
	if err := db.Save(t).Error; err != nil {
		return err
	}

	if isNew {
		column.Name = "id"
		column.Alias = "id"
		column.Type = BIGINT
		column.PrimaryKey = true
		column.TableId = t.Id
		column.NotNull = true
		column.Save()
	}
	return nil
}

func (t *Table) Delete() error {
	db := getTableDB()
	if err := db.Where("table_id = ?", t.Id).Delete(Column{}).Error; err != nil {
		return err
	}
	return db.Delete(t).Error
}

func TableDelete(where string) {
	db := getTableDB()
	db.Where(where).Delete(&Table{})
}

func TableList(tbls *[]Table) error {
	db := getTableDB()
	err := db.Find(tbls).Error
	return err
}

func (t *Table) Field(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

func (t Table) MarshalJSON() ([]byte, error) {
	j, _ := goutils.ToJsonOnly(t)
	return []byte(j), nil
}

func (t *Table) CreateTable() error {
	c := GetDriver()
	return c.CreateTable(t)
}

func (t *Table) MigrateTable() error {
	c := GetDriver()
	return c.MigrateTable(t)
}
