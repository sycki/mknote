package server

import (
	"log"
	"net/http"
	"sycki/handler"

	"github.com/gorilla/mux"
)

func StartServer() {
	m := mux.NewRouter()
	handler.BaseHandlers(m)
	handler.RestHandlers(m)
	handler.DashboardHandlers(m)
	log.Fatal(http.ListenAndServe(":80", m))
}
