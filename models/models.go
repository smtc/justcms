package models

import "github.com/smtc/JustCms/database"

func InitDB() {
	db := database.GetDB("")
	db.AutoMigrate(Account{})
	db.AutoMigrate(AccountMeta{})

	db.AutoMigrate(Post{})
	db.AutoMigrate(PostMeta{})

	db.AutoMigrate(Comment{})
	db.AutoMigrate(CommentMeta{})

	db.AutoMigrate(Link{})

	db.AutoMigrate(Options{})

	db.AutoMigrate(Term{})
	db.AutoMigrate(TermRelation{})
	db.AutoMigrate(TermTaxonomy{})
}
