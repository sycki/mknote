package handler

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

func article(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := r.Method
	if method == GET {
		tag := vars["tag"]
		en_name := vars["en_name"]

		result, err := database.GetArticle(tag, en_name)
		if err != nil {
			log.Fatal(err)
		}
		js, _ := json.Marshal(result)
		w.Write(js)
	}
}

func checkHeader(w http.ResponseWriter, r *http.Request) error {
	if debug {
		return nil
	}
	if header := r.Header.Get("x-requested-by"); header != "sycki" {
		err := errors.New(string(http.StatusNotFound))
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	return nil
}

func f(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func BaseHandlers(m *mux.Route) {
	m.HandlerFunc("/{tag:[0-9a-zA-Z-_]{1,20}}/{en_name:[0-9a-zA-Z-_]{1,50}}", f(article))
}
