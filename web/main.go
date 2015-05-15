package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/goutils"
	"github.com/smtc/justcms/auth"

	"github.com/flosch/pongo2"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	//"github.com/smtc/justcms/admin"
)

var (
	configFn = flag.String("config", "../config.json", "config file path")
)

func main() {
	config.ReadCfg(*configFn)
	deferinit.InitAll()

	//models.InitDB()
	run()
}

func run() {
	// route /admin
	//goji.Handle("/admin/*", admin.AdminMux())
	//goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	// static files
	goji.Get("/assets/*", http.FileServer(http.Dir("./")))

	goji.NotFound(goutils.NotFound)

	goji.Get("/", index)
	goji.Get("/login", login)
	goji.Post("/login", postLogin)
	goji.Get("/join", signup)
	goji.Post("/join", postSignup)

	goji.Serve()
}

/////////////////////////////////////////////////////////////////////
// just for test
func index(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl := pongo2.Must(pongo2.FromFile("./themes/default/index.html"))
	user, err := GetOrCreateUser(c, w, r)
	if err != nil {
		log.Printf("Get or Create user failed: %s\n", err.Error())
		user = auth.VisitUser
	}
	err = tpl.ExecuteWriter(pongo2.Context{
		"title": "JustCms demo",
		"user":  user,
	}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 登陆页面
func login(w http.ResponseWriter, r *http.Request) {
	tpl := pongo2.Must(pongo2.FromFile("./themes/default/login.html"))
	err := tpl.ExecuteWriter(pongo2.Context{}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 提交登陆
func postLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	msisdn := r.FormValue("msisdn")
	passwd := r.FormValue("password")

	u, err := auth.Signin(msisdn, passwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	auth.SaveUserSession(w, r, u, 0)
	http.Redirect(w, r, "/", 302)
}

// 注册页面
func signup(w http.ResponseWriter, r *http.Request) {
	tpl := pongo2.Must(pongo2.FromFile("./themes/default/signup.html"))
	err := tpl.ExecuteWriter(pongo2.Context{}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 提交注册
func postSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	msisdn := r.FormValue("msisdn")
	name := r.FormValue("username")
	passwd := r.FormValue("password")
	passwd2 := r.FormValue("password_again")

	_, err := auth.Signup(msisdn, name, passwd, passwd2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", 302)
}

/*
func get_posts() []*posts.Post {
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
*/
