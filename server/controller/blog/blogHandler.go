package blog

import (
	"mknote/server/ctx"
	"mknote/server/persistent"
	"mknote/server/persistent/structs"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var (
	model   *structs.Model
	htmlDir string
)

func init() {
	model = &structs.Model{make(map[string]interface{})}
	htmlDir = ctx.Get("html.dir")
}

func Home(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == GET {
		artID, err := persistent.GetLatestArticleID()
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, "/articles/"+artID, 302)
		}
	}
}

func Article(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == GET {
		htmlFile := htmlDir + "/article.html"
		http.ServeFile(w, r, htmlFile)
	}
}
