package models

import "github.com/smtc/JustCms/models/db"

func InitDB() {
	db := db.GetDB("")
	//db.AutoMigrate(Account{})
	println(db)
	db.CreateTable(&Account{})
}
