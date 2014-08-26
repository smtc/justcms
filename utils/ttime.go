package utils

import (
	"errors"
	"time"
)

type TTime struct {
	time.Time
	f string
}

func (t TTime) format() string {
	if t.f == "" {
		t.f = TIMEFORMAT
	}
	return t.Time.Format(t.f)
}

func (t TTime) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

func (t TTime) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("TTime.MarshalJson: year outside of range [0,9999]")
	}
	return []byte(`"` + t.format() + `"`), nil
}

func (t *TTime) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(TIMEFORMAT, string(data))
	return
}
