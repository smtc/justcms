package models

import (
	"github.com/guotie/config"
	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

// ===========================================================

type driver interface {
	HasTable(db *gorm.DB, t *Table) bool
	CreateTable(db *gorm.DB, t *Table) error
	DropTable(db *gorm.DB, t *Table) error
	MigrateTable(db *gorm.DB, t, old *Table) error
	AddColumn(db *gorm.DB, c *Column, table string) error
	DropColumn(db *gorm.DB, columns []Column, table string) error
	ChangeColumn(db *gorm.DB, c *Column, old, table string) error
	GetPage(db *gorm.DB, t *Table, where string, page, size int) (interface{}, int, error)
	SaveEntity(db *gorm.DB, t *Table, entity map[string]interface{}) error
	RemoveEntities(db *gorm.DB, tn string, ids []int64) error
}

func GetDriver() driver {
	return &mysql{}
}

// ===========================================================

type table_database int

const (
	DEFAULT_DB table_database = iota
	ACCOUNT_DB
	DYNAMIC_DB
)

func getSchema(model table_database) string {
	db := ""
	switch model {
	case DEFAULT_DB:
		db = ""
	case DYNAMIC_DB:
		db = config.GetStringDefault("dbdynamic", "")
	}
	return db
}

func GetDB(model table_database) *gorm.DB {
	return database.GetDB(getSchema(model))
}

func InitDB() {
	db := database.GetDB("")

	db.AutoMigrate(Account{})

	db.AutoMigrate(Meta{})

	db.AutoMigrate(Post{})

	db.AutoMigrate(Reply{})

	db.AutoMigrate(Link{})

	db.AutoMigrate(Options{})

	//db.AutoMigrate(Term{})
	db.AutoMigrate(TermRelation{})
	db.AutoMigrate(TermTaxonomy{})

	db.AutoMigrate(Table{})
	db.AutoMigrate(Column{})
}
