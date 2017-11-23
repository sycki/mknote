package server

import (
	"net/http"
	"sycki/log"
	"sycki/server/blog"
	"sycki/server/rest"

	"github.com/gorilla/mux"
)

func f(h http.HandlerFunc) http.HandlerFunc {
	log.Info(h)
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func StartServer() {
	m := mux.NewRouter()
	m.HandleFunc("/api/v1/article", f(rest.Article))
	m.HandleFunc("/api/v1/index", f(rest.Index))
	m.HandleFunc("/{tag:[0-9a-zA-Z-_]+}/{en_name:[0-9a-zA-Z-_]+}", f(blog.Article))
	//	m.HandleFunc("/en_name", f(blog.Article))

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(":80", nil))
}
