package handler

import (
	"net/http"
)

func dashboard(w http.ResponseWriter, r *http.Request) {

}

func DashboardHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/dashboard", f(dashboard))
}
