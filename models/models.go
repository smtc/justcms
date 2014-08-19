package models

import "github.com/smtc/JustCms/models/db"

func InitDB() {
	db := db.GetDB("")
	//db.AutoMigrate(Account{})
	db.CreateTable(Account{})
}
