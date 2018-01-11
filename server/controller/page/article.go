package page

import (
	"mknote/server/persistent"
	"mknote/server/view"
	"net/http"
)

const (
	get  = "GET"
	post = "POST"
)

func Home(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == get {
		artID, err := persistent.GetLatestArticleID()
		if err != nil {
			panic(err)
		} else {
			http.Redirect(w, r, "/articles/"+artID, 302)
		}
	}
}

func Article(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artId := r.URL.Path[len("/articles/"):]

	if method == get {
		//		htmlFile := htmlDir + "/article.html"
		//		http.ServeFile(w, r, htmlFile)
		article, err := persistent.GetArticle(artId)
		if err != nil {
			panic(err)
		}
		articleIndex, err := persistent.GetTags()
		model := &map[string]interface{}{
			"article": article,
			"index":   articleIndex,
		}
		view.RendHTML(w, "article", model)
	}
}
