package auth

import (
	"testing"

	"github.com/guotie/config"
	"github.com/smtc/justcms/database"
)

func TestMain(m *testing.M) {
	config.ReadCfg("../config.json")

	db := database.GetDB("")
	db.DropTable(&User{})
	db.AutoMigrate(&User{})

	m.Run()
}
