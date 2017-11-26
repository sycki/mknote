package rest

import (
	"encoding/json"
	"mknote/database"
	"mknote/log"
	"net/http"
	"sync"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var (
	l sync.Mutex
)

func noAuth(w http.ResponseWriter, r *http.Request) (no bool) {
	no = true
	if header := r.Header.Get("x-requested-by"); header != "mknote" {
		log.Info("api request authenticate not pass:", header, r.URL)
		http.Error(w, "authenticate not pass", http.StatusNotFound)
		return
	}
	return false
}

func Index(w http.ResponseWriter, r *http.Request) {
	if noAuth(w, r) {
		return
	}
	method := r.Method
	if method == GET {
		articleTag, err := database.GetTags()
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		result, _ := json.Marshal(articleTag)
		w.Write(result)
	}
}

func Article(w http.ResponseWriter, r *http.Request) {
	if noAuth(w, r) {
		return
	}

	method := r.Method
	uri := r.URL.Path[len("/api/v1/articles"):]
	if method == GET {
		articleTag, err := database.GetArticle(uri)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		result, _ := json.Marshal(articleTag)
		w.Write(result)
	} else if method == POST {
		articleTag, err := database.GetArticle(uri)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		result, _ := json.Marshal(articleTag)
		w.Write(result)
	}
}

func Like(w http.ResponseWriter, r *http.Request) {
	if noAuth(w, r) {
		return
	}

	method := r.Method
	artID := r.URL.Path[len("/api/v1/like"):]

	if method == PUT {
		l.Lock()
		defer l.Unlock()
		art, err := database.GetArticle(artID)
		if err != nil {
			log.Error("failed get old like number for article:", err)
			return
		}
		art.Like_count += 1
		if _, err2 := database.UpdateArtcile(art); err != nil {
			log.Error("failed write new like number for article:", err2)
		}
	}
}
