package main

import (
	"flag"
	"net/http"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/justcms/models"
	"github.com/smtc/justcms/utils"
	"github.com/zenazn/goji"

	"github.com/smtc/justcms/admin"
	_ "github.com/smtc/justcms/models"
)

var (
	configFn = flag.String("config", "./config.json", "config file path")
)

func main() {
	config.ReadCfg(*configFn)
	deferinit.InitAll()

	models.InitDB()
	run()
}

func run() {
	// route /admin
	goji.Handle("/admin/*", admin.AdminMux())
	goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	// static files
	goji.Get("/assets/*", http.FileServer(http.Dir("./")))

	goji.NotFound(utils.NotFound)

	goji.Serve()
}

func index(w http.ResponseWriter, r *http.Request) {

}
