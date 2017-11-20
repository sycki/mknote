package handler

import (
	"log"
	"net/http"
	"sycki/database"

	"github.com/gorilla/mux"
)

func indexV1(w http.ResponseWriter, r *http.Request) {
	if checkHeader(w, r) != nil {
		return
	}
	method := r.Method
	if method == GET {
		result, err := database.Index()
		if err != nil {
			log.Fatal(err)
		}
		w.Write(result)
	}
}

func articleV1(w http.ResponseWriter, r *http.Request) {
	if checkHeader(w, r) != nil {
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	method := r.Method
	if method == GET {

		result, err := database.JGET(key, ".")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(result)
	}
}

func RestHandlers(m *mux.Route) {
	m.HandlerFunc("/api/v1/article", f(articleV1))
	m.HandlerFunc("/api/v1/index", f(indexV1))
}
