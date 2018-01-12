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

package rest

import (
	"encoding/json"
	"mknote/server/persistent"
	"net/http"
	"sync"
)

const (
	get  = "GET"
	post = "POST"
	del  = "DELETE"
	put  = "PUT"
)

var (
	l sync.Mutex
)

func Index(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == get {
		articleIndex, err := persistent.GetTags()
		if err != nil {
			panic(err)
		}
		result, _ := json.Marshal(articleIndex)
		w.Write(result)
	}
}

func Article(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	uri := r.URL.Path[len("/api/v1/articles/"):]
	if method == get {
		articleTag, err := persistent.GetArticle(uri)
		if err != nil {
			panic(err)
		}
		result, _ := json.Marshal(articleTag)
		w.Write(result)
	}
}

func Like(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artID := r.URL.Path[len("/api/v1/like/"):]

	if method == post {
		l.Lock()
		defer l.Unlock()
		art, err := persistent.GetArticle(artID)
		if err != nil {
			panic(err)
		}
		art.Like_count += 1
		if _, err2 := persistent.UpdateArtcile(art); err != nil {
			panic(err2)
		}
	}
}

func Visit(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	artID := r.URL.Path[len("/api/v1/visit/articles/"):]

	if method == post {
		l.Lock()
		defer l.Unlock()
		art, err := persistent.GetArticle(artID)
		if err != nil {
			panic(err)
		}
		art.Viewer_count += 1
		if _, err2 := persistent.UpdateArtcile(art); err != nil {
			panic(err2)
		}
	}
}
