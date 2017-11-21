package blog

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sycki/database"
	"sycki/structs"

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

//func home(w http.ResponseWriter, r *http.Request) {
//	method := r.Method
//	if method == GET {
//		if !isLatestIndex {
//			result, err := database.Index()
//			if err != nil {
//				log.Fatal(err)
//			}
//			index = result
//			isLatestIndex = true
//		}
//		model.Clear()
//		model.Set("index", index)
//		view.SendHTML(w, "home")
//	}
//}

func Article(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := r.Method
	if method == GET {
		tag := vars["tag"]
		en_name := vars["en_name"]

		result, err := database.GetArticle(tag, en_name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		js, _ := json.Marshal(result)
		w.Write(js)
	}
}
