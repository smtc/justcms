package models

import (
	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

func GetDB(model string) *gorm.DB {
	db := ""
	switch model {
	case "table", "account":
		db = ""
	}
	return database.GetDB(db)
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
