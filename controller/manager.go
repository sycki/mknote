/*
Copyright 2017 sycki.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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