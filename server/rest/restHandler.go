package rest

import (
	"encoding/json"
	"mknote/database"
	"mknote/log"
	"mknote/structs"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var (
	debug bool
	model *structs.Model
)

func init() {
	debug = true
	model = &structs.Model{make(map[string]interface{})}
}

func noAuth(w http.ResponseWriter, r *http.Request) (no bool) {
	no = true
	if debug {
		return false
	}
	if header := r.Header.Get("x-requested-by"); header != "sycki" {
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

	vars := mux.Vars(r)
	key := vars["key"]
	method := r.Method
	if method == GET {
		articleTag, err := database.GetArticle(key, ".")
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		result, _ := json.Marshal(articleTag)
		w.Write(result)
	}
}
