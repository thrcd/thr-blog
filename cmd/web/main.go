package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const techPostsFS = "content/blog/tech"

const blogFS = "content/blog"

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	server := &http.Server{
		Addr:         *addr,
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("starting server", slog.String("addr", server.Addr))
	err := server.ListenAndServe()

	log.Error(err.Error())
	os.Exit(1)
}
