package utils

import (
	"bytes"
	"net/http"
)

type RequestStruct struct {
	*http.Request
	Form *GetStruct
}

func Request(r *http.Request) *RequestStruct {
	req := RequestStruct{r, Getter(r.Form)}
	return &req
}

func (r *RequestStruct) FormatBody(v interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return ToStruct(buf.Bytes(), v)
}
