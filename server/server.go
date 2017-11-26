package server

import (
	"context"
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

	//	m := mux.NewRouter()
	m := http.NewServeMux()
	m.HandleFunc("/assets/", f(static))
	m.HandleFunc("/articles/", f(blog.Article))
	m.HandleFunc("/api/v1/articles/", f(rest.Article))
	m.HandleFunc("/api/v1/like/", f(rest.Article))
	m.HandleFunc("/api/v1/index", f(rest.Index))
	//	http.Handle("/", m)

	h := &http.Server{Addr: ":80", Handler: m}

	log.Info("start mknode server...")
	go func() {
		log.Fatal(h.ListenAndServe())
	}()
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
