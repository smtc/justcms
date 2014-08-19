package utils

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/guotie/config"
)

var (
	templatePath = config.GetStringDefault("templatePath", "./assets/templates")
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".html") {
		RenderHtml(r.URL.Path, w, r)
		return
	}

	http.Error(w, "page not found.", 404)
}

func RenderHtml(path string, w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile(templatePath + path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(buf)
}
