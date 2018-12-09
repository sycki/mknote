package controller

import (
	"net/http"
)

func (m *Manager) Assets(w http.ResponseWriter, r *http.Request) {
	file := m.config.StaticDir + r.URL.Path
	http.ServeFile(w, r, file)
}

func (m *Manager) Download(w http.ResponseWriter, r *http.Request) {
	file := m.config.DownloadDir + "/" + r.URL.Path[len("/f/"):]
	http.ServeFile(w, r, file)
}
