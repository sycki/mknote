/*
Copyright 2017 sycki.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
