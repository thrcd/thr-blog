package main

import (
	"fmt"
	"github.com/thrcd/thr-blog/internal/ui"
	"net/http"
	"os"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	templateCache, err := newTemplateCache()
	if err != nil {
		fmt.Printf("got err: %s", ErrCreateTemplateCache)
		os.Exit(1)
	}

	handlers := handlers{templateCache: templateCache}

	mux.Handle("GET /static/", http.StripPrefix("/static/",
		http.FileServer(http.FS(ui.StaticFS)),
	))

	mux.HandleFunc("GET /", handlers.handleHome())
	mux.HandleFunc("GET /tech", handlers.handlePosts(blogFS))
	mux.HandleFunc("GET /life", handlers.handlePosts(blogFS))
	mux.HandleFunc("GET /posts/{type}/{dir}/{fn}", handlers.handlePost(blogFS))
	mux.HandleFunc("GET /about", handlers.handleAbout(blogFS))

	return mux
}
