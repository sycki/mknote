package blog

import (
	"mknote/config"
	"mknote/database"
	"mknote/structs"
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
	htmlDir = config.Get("html.dir")
}

func Home(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == GET {
		artID, err := database.GetLatestArticleID()
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
