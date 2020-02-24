package controller

import (
	"net/http"
)

func (m *Manager) Static(w http.ResponseWriter, r *http.Request) {
	file := m.config.HtmlDir + r.URL.Path
	http.ServeFile(w, r, file)
}

func (m *Manager) Download(w http.ResponseWriter, r *http.Request) {
	file := m.config.DownloadDir + "/" + r.URL.Path[len("/f/"):]
	http.ServeFile(w, r, file)
}
