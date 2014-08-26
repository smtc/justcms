package utils

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTTime(t *testing.T) {
	type User struct {
		Birthday TTime
	}

	tm := TTime{time.Now(), ""}

	func() {
		//s, _ := tm.MarshalJson()
		//println(string(s))
	}()

	user := User{Birthday: tm}
	s, err := json.Marshal(user)
	if err == nil {
		println(string(s))
	}
}
