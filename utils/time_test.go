package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	type User struct {
		Birthday Time
	}

	sbd := "2010-01-01 00:00:00"
	birthday, _ := time.Parse(TIMEFORMAT, sbd)

	tm := Time{birthday, ""}

	s, err := tm.MarshalJSON()
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(s) != (`"` + sbd + `"`) {
		t.Fatal(fmt.Sprintf("birthday should be %v, not %v", sbd, string(s)))
	}

	user := User{Birthday: tm}
	s, err = json.Marshal(user)
	if err != nil {
		t.Fatal(err.Error())
	}

	user2 := User{}
	if err = json.Unmarshal(s, &user2); err != nil {
		t.Fatal(err.Error())
	}
	if user2.Birthday.Time != birthday {
		t.Fatal(fmt.Sprintf("birthday should be %v, not %v", birthday, user2.Birthday))
	}

	tm.SetFormat(time.Kitchen)
	s, err = tm.MarshalText()
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(s) != "12:00AM" {
		t.Fatal(fmt.Sprintf("tm should be 12:00AM, not %v", string(s)))
	}
}
