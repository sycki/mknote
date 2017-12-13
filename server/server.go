package server

import (
	"context"
	"mknote/server/controller/blog"
	"mknote/server/controller/rest"
	"mknote/server/ctx"
	"net/http"
	"time"
)

var staticDir = ctx.Get("static.dir")

func static(w http.ResponseWriter, r *http.Request) {
	file := staticDir + r.URL.Path
	//	if  strings.Contains(file, "..") {
	//		log.Error(errInfo, file)
	//		http.NotFound(w, r)
	//		return
	//	}
	http.ServeFile(w, r, file)
}

func s(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				ctx.Error(r.RemoteAddr, "=>", r.RequestURI, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func f(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Info(r.RemoteAddr, "=>", r.RequestURI, r.UserAgent())
		defer func() {
			if err, ok := recover().(error); ok {
				ctx.Error(r.RemoteAddr, "=>", r.RequestURI, r.UserAgent(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func redirect80(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
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

	// page
	mux.HandleFunc("/", f(blog.Home))
	mux.HandleFunc("/articles/", f(blog.Article))

	// restful API
	mux.HandleFunc("/api/v1/articles/", f(rest.Article))
	mux.HandleFunc("/api/v1/index", f(rest.Index))
	mux.HandleFunc("/api/v1/like/", f(rest.Like))
	mux.HandleFunc("/api/v1/visit/", f(rest.Visit))

	// static resource
	mux.HandleFunc("/assets/", s(static))

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
	go s.ListenAndServeTLS(ctx.Get("server.tls.cert.file"), ctx.Get("server.tls.key.file"))

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
