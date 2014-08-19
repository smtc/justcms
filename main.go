package main

import (
	"flag"
	"net/http"

	"github.com/guotie/config"
	"github.com/smtc/JustCms/models"
	"github.com/smtc/JustCms/utils"
	"github.com/zenazn/goji"

	"./admin"
	_ "./models"
)

var (
	configFn = flag.String("config", "./config.json", "config file path")
)

func main() {
	config.ReadCfg(*configFn)
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
