package models

import "github.com/smtc/JustCms/database"

func init() {
	database.OpenDB("justcms_test", "root", "123456")

	//dropTables()

	InitDB()
}

func dropTables() {
	db := database.GetDB("")
	db.DropTableIfExists(Account{})
	db.DropTableIfExists(AccountMeta{})

	db.DropTableIfExists(Post{})
	db.DropTableIfExists(PostMeta{})

	db.DropTableIfExists(Reply{})
	db.DropTableIfExists(ReplyMeta{})

	db.DropTableIfExists(Link{})

	db.DropTableIfExists(Options{})

	db.DropTableIfExists(Term{})
	db.DropTableIfExists(TermRelation{})
	db.DropTableIfExists(TermTaxonomy{})
}
