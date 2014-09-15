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

func (rs *RequestStruct) FormatBody(v interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(rs.Body)
	return ToStruct(buf.Bytes(), v)
}
