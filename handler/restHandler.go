package handler

import (
	"io"
	"log"
	"net/http"
	"sycki/database"
)


func index(w http.ResponseWriter, r *http.Request) {
	if checkHeader(w, r) != nil {
		return
	}
	method := r.Method
	if method == GET {
		if !isLatestIndex {
			result, err := database.Index()
			if err != nil {
				log.Fatal(err)
			}
			indexs = result
			isLatestIndex = true
		}
		model.Clear()
		model.Set("index", indexs)
		io.WriteString(w, indexs)
	}
}

func article(w http.ResponseWriter, r *http.Request) {
	if checkHeader(w, r) != nil {
		return
	}
	key := r.FormValue("key")
	method := r.Method
	if method == GET {
		result, err := database.JGET(key, ".")
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, result)
	}
}

func RestHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/article", f(article))
	mux.HandleFunc("/index", f(index))
}
