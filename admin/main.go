package admin

import (
	"net/http"

	"../utils"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const (
	ADMIN_ROUTE = "/admin/"
)

func AdminMux() *web.Mux {
	mux := web.New()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get(ADMIN_ROUTE, indexHandler)

	mux.NotFound(utils.NotFound)
	return mux
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello world"))
}
