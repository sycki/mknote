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

package page

import (
	"html/template"
	"mknote/server/persistent"
	"mknote/server/view"
	"net/http"
	"github.com/russross/blackfriday.v2"
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
		articleHTML := template.HTML(string(blackfriday.Run([]byte(article.Content))))
		article.Content = ""

		articleIndex, err := persistent.GetTags()
		model := &map[string]interface{}{
			"articleHTML": articleHTML,
			"article":     article,
			"index":       articleIndex,
		}
		view.RendHTML(w, "article", model)
	}
}
