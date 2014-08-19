package models

import "github.com/smtc/JustCms/database"

func InitDB() {
	db := database.GetDB("")
	db.AutoMigrate(Account{})
}
