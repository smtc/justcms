package models

import (
	"testing"

	"github.com/smtc/justcms/database"
)

func TestAccount(t *testing.T) {
	account := Account{
		ObjectId: database.ObjectID(),
		Name:     "test",
		Email:    "admin@test.com",
	}

	if err := account.Save(); err != nil {
		t.Fatal(err.Error())
	}

	account2, err := AccountGet(account.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if account2.Email != account.Email {
		t.Fatal("error")
	}

	if err = account2.Delete(); err != nil {
		t.Fatal(err.Error())
	}

	account.Id = 0
	if err := account.Save(); err != nil {
		t.Fatal(err.Error())
	}
	println(account.Id)
}
