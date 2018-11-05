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

package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sycki/mknote/cmd/mknote/options"
	"github.com/sycki/mknote/controller"
	"github.com/sycki/mknote/logger"
)

// Server listen service port and map requests to handlers
type Server struct {
	mux            *http.ServeMux
	config         *options.Config
	httpServer     *http.Server
	redirectServer *http.Server
}

// NewServer initializes a new server from config
func NewServer(config *options.Config, cm *controller.Manager) *Server {
	mux := http.NewServeMux()

	// page handler
	mux.HandleFunc("/", securityHandler(cm.Index))
	mux.HandleFunc("/articles/", securityHandler(cm.Home))
	mux.HandleFunc("/f/", securityStaticHandler(cm.Download))

	// static resource
	mux.HandleFunc("/assets/", securityStaticHandler(cm.Assets))

	// restful API
	mux.HandleFunc("/v1/index", securityRest(cm.ArticleNavigation))
	mux.HandleFunc("/v1/articles/", securityRest(cm.Article))
	mux.HandleFunc("/v1/visit/articles/", securityRest(cm.Visit))
	mux.HandleFunc("/v1/like/", securityRest(cm.Like))

	mux.HandleFunc("/v1/manage/pprof/", securityRest(cm.Pprof))
	return &Server{
		config: config,
		mux:    mux,
	}
}

func securityStaticHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				debug.PrintStack()
				http.Error(w, "406", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func securityRest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"] ")
		if header := r.Header.Get("x-requested-by"); header != "mknote" {
			logger.Error(r.RemoteAddr, "=>", r.RequestURI, "unauthorized:", header, r.URL)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				debug.PrintStack()
				http.Error(w, "500", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func securityHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Method, r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"] ")
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				debug.PrintStack()
				http.Error(w, "500", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func (s *Server) redirect80(w http.ResponseWriter, r *http.Request) {
	hostname := s.config.HostName
	if hostname == "" {
		hostname = r.Host
	}
	target := "https://" + hostname + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

// Start the server
func (s *Server) Start(errCh chan error) {
	if s.config.IsTls {
		logger.Info(fmt.Sprintf("starting server on port %s and %s ...", s.config.TlsPort, s.config.HttpPort))
		s.httpServer = &http.Server{Addr: ":" + s.config.TlsPort, Handler: s.mux, ErrorLog: logger.GetLogger()}
		go s.httpServer.ListenAndServeTLS(s.config.TlsCertFile, s.config.TlsKeyFile)

		logger.Info("redirect all request to tls from http")
		s.redirectServer = &http.Server{Addr: ":" + s.config.HttpPort, Handler: http.HandlerFunc(s.redirect80), ErrorLog: logger.GetLogger()}
		go s.redirectServer.ListenAndServe()
	} else {
		logger.Info(fmt.Sprintf("starting server on port %s ...", s.config.HttpPort))
		s.httpServer = &http.Server{Addr: ":" + s.config.HttpPort, Handler: s.mux, ErrorLog: logger.GetLogger()}
		go s.httpServer.ListenAndServe()
	}
}

// Stop the server
func (s *Server) Stop() {
	logger.Info("stopping server ...")
	if s.httpServer != nil {
		s.httpServer.Shutdown(context.Background())
	}
	if s.config.IsTls && s.redirectServer != nil {
		s.redirectServer.Shutdown(context.Background())
	}
}
