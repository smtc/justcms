package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/goutils"
	"github.com/smtc/justcms/models"
	"github.com/zenazn/goji"

	"github.com/flosch/pongo2"
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

	goji.NotFound(goutils.NotFound)

	goji.Get("/", index)

	goji.Serve()
}

// just for test
func index(w http.ResponseWriter, r *http.Request) {
	tpl := pongo2.Must(pongo2.FromFile("./themes/default/index.html"))
	err := tpl.ExecuteWriter(pongo2.Context{
		"title":     "JustCms demo",
		"get_posts": get_posts,
		"show_post": show_post,
	}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get_posts() []*models.Post {
	return []*models.Post{
		&models.Post{
			Id:      1,
			Title:   "a test post",
			Content: "JustCms is written by golang, it's awesome",
		},
		&models.Post{
			Id:      2,
			Title:   "冷核聚变技术取得总大突破！！！",
			Content: "历史性的一刻，人类再也不用为能源发愁了。",
		},
	}
}

func show_post(post *models.Post) string {
	return fmt.Sprintf(`<h1>%s</h1>
	<p>%s</p>
	`, post.Title, post.Content)
}
