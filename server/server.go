package server

import (
	"context"
	"mknote/config"
	"mknote/log"
	"mknote/server/blog"
	"mknote/server/rest"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	assDir = "static"
)

func static(w http.ResponseWriter, r *http.Request) {
	file := assDir + r.URL.Path
	//	if  strings.Contains(file, "..") {
	//		log.Error(errInfo, file)
	//		http.NotFound(w, r)
	//		return
	//	}
	http.ServeFile(w, r, file)
}

func f(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func StartServer() {
	log.Info("load handlers...")
	m := http.NewServeMux()

	// page
	m.HandleFunc("/", f(blog.Home))
	m.HandleFunc("/articles/", f(blog.Article))

	// restful API
	m.HandleFunc("/api/v1/articles/", f(rest.Article))
	m.HandleFunc("/api/v1/like/", f(rest.Like))
	m.HandleFunc("/api/v1/index", f(rest.Index))

	// static resource
	m.HandleFunc("/assets/", f(static))

	var isTls = true
	tlsCert := config.Get("server.tls.cert.file")
	tlsKey := config.Get("server.tls.key.file")

	// test tls files is exists.
	f1, e1 := os.OpenFile(tlsCert, os.O_RDONLY, 0666)
	if e1 != nil {
		log.Warn("failed try open cert file", tlsCert, e1)
		isTls = false
	} else {
		f1.Close()
	}
	f2, e2 := os.OpenFile(tlsKey, os.O_RDONLY, 0666)
	if e2 != nil {
		log.Warn("failed try open key file", tlsKey, e2)
		isTls = false
	} else {
		f2.Close()
	}

	var h *http.Server

	if isTls {
		log.Info("start server with https")
		h = &http.Server{Addr: ":443", Handler: m}
		go func() {
			log.Fatal(h.ListenAndServeTLS(tlsCert, tlsKey))
		}()
	} else {
		go func() {
			log.Info("start server with http")
			h = &http.Server{Addr: ":80", Handler: m}
			log.Fatal(h.ListenAndServe())
		}()
	}
	log.Info("mknode started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Warn("recived signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := h.Shutdown(ctx)
	if err != nil {
		log.Warn(err)
	}

	log.Info("mknode gracefully stopped")
	println("mknode gracefully stopped")
}
