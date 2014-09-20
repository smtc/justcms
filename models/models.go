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
	db.AutoMigrate(AccountMeta{})

	db.AutoMigrate(Post{})
	db.AutoMigrate(PostMeta{})

	db.AutoMigrate(Reply{})
	db.AutoMigrate(ReplyMeta{})

	db.AutoMigrate(Link{})

	db.AutoMigrate(Options{})

	//db.AutoMigrate(Term{})
	db.AutoMigrate(TermRelation{})
	db.AutoMigrate(TermTaxonomy{})

	db.AutoMigrate(Table{})
	db.AutoMigrate(Column{})
}
