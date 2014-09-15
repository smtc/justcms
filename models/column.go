package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

type Column struct {
	Id        int64     `json:"id"`
	TableId   int64     `sql:"not null" json:"table_id"`
	Name      string    `sql:"size:45;not null" json:"name"`
	Alias     string    `sql:"size:45" json:"alias"`
	Des       string    `sql:"size:512" json:"des"`
	Type      string    `sql:"size:20" json:"type"`
	Size      int       `json:"size"`
	Filter    string    `sql:"size:127" json:"filter"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
}

func getColumnDB() *gorm.DB {
	return database.GetDB(dynamic_db)
}

func (c *Column) Exist() bool {
	var count int
	db := getColumnDB()
	db.Model(Column{}).Where("id!=? and table_id=? and name=?", c.Id, c.TableId, c.Name).Count(&count)
	return count > 0
}

func (c *Column) Get(id int64) error {
	db := getColumnDB()
	return db.First(c, id).Error
}

func (c *Column) Save() error {
	db := getColumnDB()
	if c.Exist() {
		return fmt.Errorf("column '%v' is existed", c.Name)
	}
	c.EditAt = time.Now()
	return db.Save(c).Error
}

func (c *Column) Delete() error {
	db := getColumnDB()
	return db.Delete(c).Error
}

func ColumnDelete(where string) {
	db := getColumnDB()
	db.Where(where).Delete(&Column{})
}

func ColumnList(tableId int64) ([]Column, error) {
	db := getColumnDB()
	var cols []Column

	err := db.Where("table_id=?", tableId).Find(&cols).Error
	return cols, err
}
