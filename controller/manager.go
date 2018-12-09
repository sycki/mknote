package controller

import (
	"context"
	"github.com/sycki/mknote/logger"
	"net/http"
	"sync"

	"github.com/sycki/mknote/cmd/mknote/options"
	"github.com/sycki/mknote/storage"
	"github.com/sycki/mknote/view"
)

const (
	get  = "GET"
	post = "POST"
	del  = "DELETE"
	put  = "PUT"
)

type Manager struct {
	l       sync.Mutex
	view    *view.View
	config  *options.Config
	storage *storage.Manager
	pprof   *http.Server
}

func NewManager(conf *options.Config, storage *storage.Manager) *Manager {
	return &Manager{
		config:  conf,
		view:    view.NewView(conf),
		storage: storage,
	}
}

func (m *Manager) Close() {
	if m.pprof != nil {
		m.pprof.Shutdown(context.Background())
		logger.Info("pprofile server is stopped gracefully")
	}
}