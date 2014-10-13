package models

import (
	//"encoding/json"
	//"fmt"
	"github.com/smtc/justcms/database"
	"testing"
)

//
func TestScan(t *testing.T) {
	u1 := Account{Name: "guotie",
		ObjectId: database.ObjectID(),
		Capability: &AccountCap{
			Roles: []string{"admin", "manager"},
			Caps: map[string]bool{
				"read":  true,
				"write": false,
			},
		},
	}
	err := u1.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = u1.SetCaps()
	if err != nil {
		t.Fatal(err)
	}

	var u2 Account
	if err = u2.Get(u1.Id); err != nil {
		t.Fatal(err)
	}
	if err = u2.GetCaps(); err != nil {
		t.Fatal(err)
	}
}
