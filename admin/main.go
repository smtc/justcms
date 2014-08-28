package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/smtc/justcms/utils"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func AdminMux() *web.Mux {
	mux := web.New()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/admin/", indexHandler)
	mux.Get("/admin/menu", menuHandler)

	mux.Get("/admin/account/", AccountList)
	mux.Get("/admin/account/:id", AccountEntity)

	mux.Get(regexp.MustCompile(`^/admin/(?P<model>.+)\.(?P<fn>.+)$`), templateHandler)

	mux.NotFound(utils.NotFound)
	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderHtml("/admin/main.html", w, r)
}

/*
模板页暂时以 model.fn 分级
*/
func templateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	temp := fmt.Sprintf("/admin/%s_%s.html", c.URLParams["model"], c.URLParams["fn"])
	utils.RenderHtml(temp, w, r)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("./admin/menu.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
