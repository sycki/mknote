package handler

import (
	"errors"
	"net/http"
	"sycki/structs"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var (
	debug         bool
	model         *structs.Model
	indexs        string
	isLatestIndex bool
)

func init() {
	debug = true
	model = &structs.Model{make(map[string]interface{})}
	isLatestIndex = false
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

func BaseHandlers(mux *http.ServeMux) {
	//	mux.HandleFunc("/article", f(article))
}
