package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func dashboard(w http.ResponseWriter, r *http.Request) {

}

func DashboardHandlers(m *mux.Route) {
	m.HandlerFunc("/dashboard", f(dashboard))
}
