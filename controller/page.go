package controller

import (
	"html/template"
	"net/http"
	"github.com/russross/blackfriday.v2"
)

func (m *Manager) Index(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == get {
		artID, err := m.storage.GetLatestArticleID()
		if err != nil {
			panic(err)
		} else {
			http.Redirect(w, r, "/articles/"+artID, 302)
		}
	}
}

func (m *Manager) Home(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artId := r.URL.Path[len("/articles/"):]

	if method == get {
		article, err := m.storage.GetArticle(artId)
		if err != nil {
			panic(err)
		}
		articleHTML := template.HTML(string(blackfriday.Run([]byte(article.Content))))
		article.Content = ""

		articleIndex, err := m.storage.GetTags()
		model := &map[string]interface{}{
			"articleHTML": articleHTML,
			"article":     article,
			"index":       articleIndex,
		}
		m.view.RendHTML(w, "article", model)
	}
}
