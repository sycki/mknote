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
	"github.com/sycki/mknote/server/controller/page"
	"github.com/sycki/mknote/server/controller/rest"
	"github.com/sycki/mknote/logger"
	"net/http"
	"github.com/sycki/mknote/cmd/mknote/options"
)

var (
	config *options.Config
)

func securityStaticHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(r.RemoteAddr, "=>", r.RequestURI, err)
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
				logger.Error(r.RemoteAddr, "=>", r.RequestURI, err)
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
				logger.Error(r.Method, r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"]", err)
				http.Error(w, "500", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func redirect80(w http.ResponseWriter, r *http.Request) {
	hostname := config.HostName
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

func buildMux() *http.ServeMux {
	logger.Info("load handlers...")
	mux := http.NewServeMux()

	// page handler
	mux.HandleFunc("/", securityHandler(page.Home))
	mux.HandleFunc("/articles/", securityHandler(page.Article))
	mux.HandleFunc("/uploads/", securityStaticHandler(rest.Download))

	// restful API
	mux.HandleFunc("/api/v1/index", securityRest(rest.Index))
	mux.HandleFunc("/api/v1/articles/", securityRest(rest.Article))
	mux.HandleFunc("/api/v1/visit/articles/", securityRest(rest.Visit))
	mux.HandleFunc("/api/v1/like/", securityRest(rest.Like))

	// static resource
	mux.HandleFunc("/assets/", securityStaticHandler(rest.Assets))

	return mux
}

func Start(conf *options.Config, c context.Context) {
	config = conf
	logger.Info("start server with http")

	s := &http.Server{Addr: ":80", Handler: buildMux(), ErrorLog: logger.GetLogger()}
	go s.ListenAndServe()
	logger.Info("mknode started")

	<-c.Done()
	s.Shutdown(c)
}

func StartTLS(conf *options.Config, c context.Context) {
	config = conf
	logger.Info("start server with tls")

	s := &http.Server{Addr: ":443", Handler: buildMux(), ErrorLog: logger.GetLogger()}
	go s.ListenAndServeTLS(config.TlsCertFile, config.TlsKeyFile)

	logger.Info("redirect 80 to 443")
	s80 := &http.Server{Addr: ":80", Handler: http.HandlerFunc(redirect80), ErrorLog: logger.GetLogger()}
	go s80.ListenAndServe()
	logger.Info("mknode started")

	<-c.Done()
	s.Shutdown(c)
	s80.Shutdown(c)
}
