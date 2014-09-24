package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

//===========================================================

const (
	AUTO_INCREMENT = "auto_increment"
	INT            = "int"
	BIGINT         = "bigint"
	FLOAT          = "float"
	DOUBLE         = "double"
	VARCHAR        = "varchar"
	TEXT           = "text"
	LONGTEXT       = "longtext"
	BOOL           = "bool"
	DATE           = "date"
	DATETIME       = "datetime"
	PICTURE        = "picture"
	PASSWORD       = "password"

	PRIVATE   = 0
	PROTECTED = 1
	PUBLIC    = 2
)

type columnType struct {
	Name string
	Size int
	Des  string
}

var ColumnTypes = map[string]columnType{
	AUTO_INCREMENT: columnType{AUTO_INCREMENT, 20, ""},
	INT:            columnType{INT, 11, ""},
	BIGINT:         columnType{BIGINT, 20, ""},
	FLOAT:          columnType{FLOAT, 0, ""},
	DOUBLE:         columnType{DOUBLE, 0, ""},
	VARCHAR:        columnType{VARCHAR, 45, ""},
	TEXT:           columnType{TEXT, 0, ""},
	LONGTEXT:       columnType{LONGTEXT, 0, ""},
	BOOL:           columnType{BOOL, 1, ""},
	DATE:           columnType{DATE, 0, ""},
	DATETIME:       columnType{DATETIME, 0, ""},
	PICTURE:        columnType{PICTURE, 255, ""},
	PASSWORD:       columnType{PASSWORD, 128, ""},
}

var Filters = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"gt":   ">",
	"egt":  ">=",
	"lt":   "<",
	"elt":  "<=",
	"like": "like",
	"in":   "in",
}

// ================================================================

type Column struct {
	Id           int64     `json:"id"`
	TableId      int64     `sql:"not null" json:"table_id"`
	Name         string    `sql:"size:45;not null" json:"name"`
	Alias        string    `sql:"size:45" json:"alias"`
	Des          string    `sql:"size:512" json:"des"`
	Type         string    `sql:"size:20" json:"type"`
	Size         int       `json:"size"`
	Filter       string    `sql:"size:127" json:"filter"`
	NotNull      bool      `json:"not_null"`
	PrimaryKey   bool      `json:"primary_key"`
	DefaultValue string    `sql:"size:512" json:"default_value"`
	CreatedAt    time.Time `json:"created_at"`
	EditAt       time.Time `json:"edit_at"`
	EditAble     bool      `json:"edit_able"`
	OrderIndex   int       `json:"order_index"`
}

func getColumnDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
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
	var (
		db    = getColumnDB()
		ndb   = GetDB(DYNAMIC_DB)
		d     = GetDriver()
		table *Table
		old   Column
		err   error
	)
	if c.Exist() {
		return fmt.Errorf("column '%v' is existed", c.Name)
	}

	table, err = GetTable(c.TableId)
	if err != nil {
		return err
	}
	if c.Id == 0 {
		c.CreatedAt = time.Now()
		if d.HasTable(ndb, table) {
			d.AddColumn(ndb, c, table.Name)
		}
	} else {
		old.Get(c.Id)

		if old.Name == "id" {
			return fmt.Errorf("column 'id' can't rename.")
		}

		c.CreatedAt = old.CreatedAt
		d.ChangeColumn(ndb, c, old.Name, table.Name)
	}

	c.EditAt = time.Now()
	err = db.Save(c).Error
	table.Refresh()
	return err
}

func (c *Column) Delete() error {
	return ColumnDelete("id = ?", c.Id)
}

func ColumnDelete(where string, data ...interface{}) error {
	var (
		db      = getColumnDB()
		ndb     = GetDB(DYNAMIC_DB)
		d       = GetDriver()
		columns []Column
		table   *Table
		err     error
	)

	err = db.Where(where, data...).Find(&columns).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, c := range columns {
		if c.Name == "id" {
			return fmt.Errorf("column 'id' can't remove.")
		}
	}

	table, _ = GetTable(columns[0].TableId)
	d.DropColumn(ndb, columns, table.Name)
	err = db.Where(where, data...).Delete(&Column{}).Error
	table.Refresh()

	return err
}

func ColumnList(tableIds []int64) ([]Column, error) {
	db := getColumnDB()
	var (
		cols  []Column
		where string
	)

	if len(tableIds) == 1 {
		where = "table_id = ?"
	} else {
		where = "table_id in (?)"
	}

	err := db.Where(where, tableIds).Order("order_index").Find(&cols).Error
	return cols, err
}

// sort =========================================================

type ColumnSort []Column

func (c ColumnSort) Len() int {
	return len(c)
}

func (c ColumnSort) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ColumnSort) Less(i, j int) bool {
	return c[i].OrderIndex < c[j].OrderIndex
}
