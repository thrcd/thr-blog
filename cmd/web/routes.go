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

	mux.HandleFunc("GET /", handlers.handlePosts(techPostsFS))
	mux.HandleFunc("GET /post/{dir}/{fn}", handlers.handlePost(techPostsFS))
	mux.HandleFunc("GET /about", handlers.handleAbout(blogFS))

	return mux
}
