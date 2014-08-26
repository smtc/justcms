package utils

import (
	"errors"
	"time"
)

type TTime struct {
	time.Time
}

func (t TTime) MarshalJson() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("TTime.MarshalJson: year outside of range [0,9999]")
	}
	return []byte(t.Format("2006-01-02 15:04:05")), nil
}

func (t *TTime) UnmarshalJson(data []byte) (err error) {
	*t, err = time.Parse("2006-01-02 15:04:05", string(data))
	return
}
