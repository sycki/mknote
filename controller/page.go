package controller

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (m *Manager) Index(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artId := "mknote/article-map"

	if method == get {
		article, err := m.storage.GetArticle(artId)
		if err != nil {
			panic(err)
		}
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.NoEmptyLineBeforeBlock | parser.HardLineBreak)
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

	if method == get {
		artId := r.URL.Path[len("/articles/"):]
		artIdLower := strings.ToLower(artId)
		if strings.HasSuffix(artIdLower, ".png") || strings.HasSuffix(artIdLower, ".jpg") || strings.HasSuffix(artIdLower, ".jpeg") || strings.HasSuffix(artIdLower, ".gif") {
			data, err := m.storage.GetMedia(artId)
			if err != nil {
				panic(err)
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000") // 长时间缓存
			w.Write(data)
			return
		}
		article, err := m.storage.GetArticle(artId)
		if err != nil {
			panic(err)
		}

		// 定义一个钩子函数，用于把文章中的<img />用<a />包裹起来。
		var inLink bool
		ReplaceImg := func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			switch nodet := node.(type) {
			case *ast.Image:
			case *ast.HTMLSpan:
				var src string
				words := strings.Fields(string(nodet.Literal))
				if len(words) < 1 {
					return ast.GoToNext, false
				}
				if words[0] == "<a" {
					inLink = true
					return ast.GoToNext, false
				}
				if words[0] == "</a>" {
					inLink = false
					return ast.GoToNext, false
				}
				if words[0] == "<img" {
					if inLink {
						return ast.GoToNext, false
					}
					for _, word := range words[1:] {
						kv := strings.Split(word, "=")
						if len(kv) < 2 {
							continue
						}
						if strings.TrimSpace(kv[0]) == "src" {
							src = strings.TrimPrefix(strings.TrimSpace(kv[1]), `"`)
							src = strings.TrimSuffix(src, `"`)
							break
						}
					}
				}
				if len(src) < 1 {
					return ast.GoToNext, false
				}
				line := fmt.Sprintf(`<a href="%s" data-lightbox="gallery">%s</a>`, src, nodet.Literal)
				nodet.Literal = []byte(line)
			}
			return ast.GoToNext, false
		}

		opts := html.RendererOptions{
			Flags:          html.CommonFlags,
			RenderNodeHook: ReplaceImg,
		}
		renderer := html.NewRenderer(opts)
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.NoEmptyLineBeforeBlock)
		articleHTML := template.HTML(markdown.ToHTML([]byte(article.Content), p, renderer))
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
