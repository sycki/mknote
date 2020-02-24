package controller

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"net/http"
)

func (m *Manager) Index(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artId := "mknote/article-map"

	if method == get {
		article, err := m.storage.GetArticle(artId)
		if err != nil {
			panic(err)
		}
		p := parser.NewWithExtensions(parser.CommonExtensions|parser.NoEmptyLineBeforeBlock)
		articleHTML := template.HTML(markdown.ToHTML([]byte(article.Content), p, nil))
		article.Content = ""

		articleIndex, err := m.storage.GetTags()
		model := &map[string]interface{}{
			"articleHTML": articleHTML,
			"article":     article,
			"index":       articleIndex,
		}
		m.view.RendHTML(w, "article.html", model)
	}
}

func (m *Manager) Article(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artId := r.URL.Path[len("/articles/"):]

	if method == get {
		article, err := m.storage.GetArticle(artId)
		if err != nil {
			panic(err)
		}
		p := parser.NewWithExtensions(parser.CommonExtensions|parser.NoEmptyLineBeforeBlock)
		articleHTML := template.HTML(markdown.ToHTML([]byte(article.Content), p, nil))
		article.Content = ""

		articleIndex, err := m.storage.GetTags()
		model := &map[string]interface{}{
			"articleHTML": articleHTML,
			"article":     article,
			"index":       articleIndex,
		}
		m.view.RendHTML(w, "article.html", model)
	}
}
