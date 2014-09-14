package utils

import (
	"bytes"
	"net/http"
)

func RequestStruct(r *http.Request, v interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return ToStruct(buf.Bytes(), v)
}
