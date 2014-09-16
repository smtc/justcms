package utils

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type httpHandler struct {
	Env *GetStruct
	P   *GetStruct
	W   *RenderStruct
	R   *RequestStruct
}

func HttpHandler(c web.C, w http.ResponseWriter, r *http.Request) *httpHandler {
	return &httpHandler{
		Env: Getter(c.Env),
		P:   Getter(c.URLParams),
		W:   Render(w),
		R:   Request(r),
	}
}
