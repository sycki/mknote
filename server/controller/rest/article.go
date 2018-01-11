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
