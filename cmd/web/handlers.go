package main

import (
	"errors"
	"net/http"
	"text/template"
)

var (
	ErrIBrokeSomething = errors.New("sorry, I probably broke something")
)

type handlers struct {
	templateCache map[string]*template.Template
}

func (h *handlers) handlePosts(postsDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dirs := getSubDirs(postsDir)
		if dirs == nil {
			h.empty(w)
			return
		}

		postListItems := getPostListItems(dirs)

		// DirNames represent the subsections on the posts page.
		// In this case, 2024, 2025...
		dirNames := maps(dirs, func(item string) string { return lastSubString(item, "/") })

		data := newTemplateData()
		data.Dirs = dirNames
		data.PostListItems = postListItems

		render(w, "posts.tmpl", h.templateCache, data)
	}
}

func (h *handlers) handlePost(postDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fn := r.PathValue("fn")
		currentDir := r.PathValue("dir")
		filePath := postDir + "/" + currentDir + "/" + fn + ".md"

		md, err := getMarkdown(filePath)
		if err != nil {
			h.serverError(w)
			return
		}

		data := newTemplateData()
		data.Markdown = md

		render(w, "post.tmpl", h.templateCache, data)
	}
}

func (h *handlers) handleAbout(aboutDir string) http.HandlerFunc {
	filePath := aboutDir + "/" + "about.md"

	return func(w http.ResponseWriter, r *http.Request) {
		md, err := getMarkdown(filePath)
		if err != nil {
			h.serverError(w)
			return
		}

		data := newTemplateData()
		data.Markdown = md
		render(w, "about.tmpl", h.templateCache, data)
	}
}

func (h *handlers) serverError(w http.ResponseWriter) {
	data := newTemplateData()
	data.Error = ErrIBrokeSomething.Error()
	renderError(w, "error.tmpl", h.templateCache, data)
}

func (h *handlers) empty(w http.ResponseWriter) {
	data := newTemplateData()
	render(w, "empty.tmpl", h.templateCache, data)
}
