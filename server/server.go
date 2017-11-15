package server

import (
	"log"
	"net/http"
	"sycki/handler"
)

func StartServer() {
	mux := http.NewServeMux()
	handler.BaseHandlers(mux)
	handler.RestHandlers(mux)
	handler.DashboardHandlers(mux)
	log.Fatal(http.ListenAndServe(":80", mux))
}
