package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/thrcd/thr-blog/internal/parser"
	"github.com/thrcd/thr-blog/internal/ui"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	ErrCreateTemplateCache = errors.New("can not create a new template cache")
)

type postListItem struct {
	Filename string
	Metadata parser.Metadata
}

type templateData struct {
	Dirs          []string
	PostListItems map[string][]postListItem
	Markdown      parser.Markdown
	CurrentDate   time.Time
	Error         string
	PostType      string
}

func newTemplateData() templateData {
	return templateData{
		CurrentDate: time.Now(),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// Rendering

func render(w http.ResponseWriter, page string, template map[string]*template.Template, data templateData) {
	ts, ok := template[page]

	if !ok {
		fmt.Errorf("the template %s does not exist", page)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return
	}

	buf.WriteTo(w)
}

func renderError(w http.ResponseWriter, page string, template map[string]*template.Template, data templateData) {
	ts, ok := template[page]

	if !ok {
		fmt.Errorf("the template %s does not exist", page)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(w, "error", data)
	if err != nil {
		return
	}

	_, _ = buf.WriteTo(w)
}

// Template Functions

func formatDate(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(format)
}

func year(t time.Time) string {
	return strconv.Itoa(t.Year())
}

func cleanMarkdownExt(fn string) string {
	return strings.TrimSuffix(fn, ".md")
}

func lowcase(str string) string {
	return strings.ToLower(str)
}

var functions = template.FuncMap{
	"formatDate":       formatDate,
	"cleanMarkdownExt": cleanMarkdownExt,
	"year":             year,
	"lowcase":          lowcase,
}
