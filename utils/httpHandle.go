package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".html") {
		http.Error(w, "page not found.", 404)
		return
	}

	templatePath := "./assets/templates"
	buf, err := ioutil.ReadFile(templatePath + r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(buf)
}
