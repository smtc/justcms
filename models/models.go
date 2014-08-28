package models

import "github.com/smtc/JustCms/database"

var (
	account_db = ""
	post_db    = ""
)

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

	db.AutoMigrate(Term{})
	db.AutoMigrate(TermRelation{})
	db.AutoMigrate(TermTaxonomy{})
}
