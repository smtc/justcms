package models

import (
	"fmt"
	"strings"
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
	db := getColumnDB()
	if c.Exist() {
		return fmt.Errorf("column '%v' is existed", c.Name)
	}
	if c.Id == 0 {
		c.CreatedAt = time.Now()
	}
	c.EditAt = time.Now()
	return db.Save(c).Error
}

func (c *Column) Delete() error {
	db := getColumnDB()
	return db.Delete(c).Error
}

func ColumnDelete(where string, data interface{}) error {
	db := getColumnDB()
	return db.Where(where, data).Delete(&Column{}).Error
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

	err := db.Where(where, tableIds).Find(&cols).Error
	return cols, err
}

func (c *Column) GetCreate() string {
	var (
		size int
		sql  = ""
		ct   = ColumnTypes[c.Type]
	)
	sql += fmt.Sprintf("`%v` ", c.Name)
	switch c.Type {
	case AUTO_INCREMENT:
		sql += "BIGINT(20) NOT NULL AUTO_INCREMENT"
		return sql
	case BOOL:
		sql += "TINYINT(1) "
	case PICTURE, VARCHAR:
		size = c.Size
		if size == 0 {
			size = ct.Size
		}
		sql += fmt.Sprintf("VARCHAR(%v) ", size)
	default:
		sql += fmt.Sprintf("%v ", strings.ToUpper(c.Type))
	}
	if c.NotNull {
		sql += "NOT NULL "
	} else {
		sql += "NULL DEFAULT NULL "
	}
	return sql
}
