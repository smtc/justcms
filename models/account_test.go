package models

import (
	"testing"
	"time"

	"github.com/smtc/justcms/database"
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

	account.Id = 0
	account.Name = "名字"
	if err := account.Save(); err != nil {
		t.Fatal(err.Error())
	}

	account2.Delete()

	if account.Name == account2.Name {
		t.Fatal(account.Name + "==" + account2.Name)
	}

	AccountDelete("name='名字'")

}
