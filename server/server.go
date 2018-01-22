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
	"mknote/server/controller/page"
	"mknote/server/controller/rest"
	"mknote/server/ctx"
	"net/http"
	"time"
)

var (
	staticDir  = ctx.Config.StaticDir
	uploadsDir = ctx.Config.UploadsDir
)

func upload(w http.ResponseWriter, r *http.Request) {
	file := uploadsDir + "/" + r.URL.Path[len("/uploads/"):]
	http.ServeFile(w, r, file)
}

func static(w http.ResponseWriter, r *http.Request) {
	file := staticDir + r.URL.Path
	//	if  strings.Contains(file, "..") {
	//		log.Error(errInfo, file)
	//		http.NotFound(w, r)
	//		return
	//	}
	http.ServeFile(w, r, file)
}

func securityStaticHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Error(r.RemoteAddr, "=>", r.RequestURI, err)
				http.Error(w, "406", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func securityRest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Info(r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"] ")
		if header := r.Header.Get("x-requested-by"); header != "mknote" {
			ctx.Error(r.RemoteAddr, "=>", r.RequestURI, "unauthorized:", header, r.URL)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		defer func() {
			if err := recover(); err != nil {
				ctx.Error(r.RemoteAddr, "=>", r.RequestURI, err)
				http.Error(w, "500", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func securityHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Info(r.Method, r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"] ")
		defer func() {
			if err := recover(); err != nil {
				ctx.Error(r.Method, r.RemoteAddr, "=>", r.RequestURI, "["+r.UserAgent()+"]", err)
				http.Error(w, "500", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func redirect80(w http.ResponseWriter, r *http.Request) {
	hostname := ctx.Config.HostName
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
	ctx.Info("load handlers...")
	mux := http.NewServeMux()

	// page handler
	mux.HandleFunc("/", securityHandler(page.Home))
	mux.HandleFunc("/articles/", securityHandler(page.Article))
	mux.HandleFunc("/uploads/", securityStaticHandler(upload))

	// restful API
	mux.HandleFunc("/api/v1/index", securityRest(rest.Index))
	mux.HandleFunc("/api/v1/articles/", securityRest(rest.Article))
	mux.HandleFunc("/api/v1/visit/articles/", securityRest(rest.Visit))
	mux.HandleFunc("/api/v1/like/", securityRest(rest.Like))

	// static resource
	mux.HandleFunc("/assets/", securityStaticHandler(static))

	return mux
}

func Start() {
	ctx.Info("start server with http")

	s := &http.Server{Addr: ":80", Handler: buildMux(), ErrorLog: ctx.GetLogger()}
	go s.ListenAndServe()

	ctx.Info("mknode started")

	stop := make(chan string, 1)
	ctx.RegistryStoper(stop)

	<-stop
	close(stop)
	ctx.Info("server read stop")

	c, _ := context.WithTimeout(context.Background(), 1*time.Second)
	s.Shutdown(c)

	//wait other hook execut complete
	time.Sleep(1 * time.Second)

	ctx.Info("server gracefully stopped")
}

func StartTLS() {
	ctx.Info("start server with tls")

	s := &http.Server{Addr: ":443", Handler: buildMux(), ErrorLog: ctx.GetLogger()}
	go s.ListenAndServeTLS(ctx.Config.TlsCertFile, ctx.Config.TlsKeyFile)

	ctx.Info("redirect 80 to 443")
	s80 := &http.Server{Addr: ":80", Handler: http.HandlerFunc(redirect80), ErrorLog: ctx.GetLogger()}
	go s80.ListenAndServe()

	ctx.Info("mknode started")

	stop := make(chan string, 1)
	ctx.RegistryStoper(stop)

	<-stop
	ctx.Info("serverTLS read stop")

	close(stop)
	c, _ := context.WithTimeout(context.Background(), 1*time.Second)
	s.Shutdown(c)
	s80.Shutdown(c)

	//wait other hook execut complete
	time.Sleep(1 * time.Second)

	ctx.Info("server gracefully stopped")
}
