package main

import (
	"net/http"

	"github.com/zenazn/goji"

	"./admin"
	_ "./models"
)

var (
	templatePath = "./assets/templates/"
)

func main() {
	run()
}

func run() {
	// route /admin
	goji.Handle("/admin/*", admin.AdminMux())
	goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	goji.Serve()
}
