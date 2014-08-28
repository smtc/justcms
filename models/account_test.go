package models

import (
	"testing"
	"time"

	"github.com/smtc/justcms/database"
	"github.com/smtc/justcms/utils"
)

func TestAccount(t *testing.T) {
	account := Account{
		ObjectId:  database.ObjectID(),
		Name:      "test",
		Email:     "admin@test.com",
		Avatar:    "admin@test.com",
		Msisdn:    "",
		Password:  "123456",
		City:      "Nanjing",
		CreatedAt: time.Now(),
	}

	if err := account.Save(); err != nil {
		t.Fatal(err.Error())
	}

	var account2 Account
	if err := account2.Get(account.Id); err != nil {
		t.Fatal(err.Error())
	}
	if account2.Email != account.Email {
		t.Fatal("error")
	}

	account.Name = "名字"
	if err := account.Save(); err != nil {
		t.Fatal(err.Error())
	}

	if account.Name == account2.Name {
		t.Fatal(account.Name + "==" + account2.Name)
	}

	s, _ := utils.ToJsonOnly(account)
	println(s)
}
