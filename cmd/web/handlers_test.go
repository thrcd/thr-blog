package main

import (
	"github.com/thrcd/thr-blog/internal/testkit"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePosts(t *testing.T) {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test handling post list request")
	{
		t.Log("Test 0: When calling posts and getting them from posts folder.")
		{
			// 2024
			want := func() string {
				dirs := getSubDirs("content/test/posts")
				return lastSubString(dirs[0], "/")
			}

			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handlePosts("content/test/posts").ServeHTTP(w, r)

			testkit.Check(t, strings.Contains(w.Body.String(), want()), "Should see page section %s in body response", want())
		}

		t.Log("Test 1: When calling post list and simulating empty posts.")
		{
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handlePosts("").ServeHTTP(w, r)

			testkit.Check(t, strings.Contains(w.Body.String(), "empty"), "Should see home content [empty] in body response.")
		}

		t.Log("Test 2: When calling post list and simulating empty directory.")
		{
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handlePosts("test/empty").ServeHTTP(w, r)

			testkit.Check(t, strings.Contains(w.Body.String(), "empty"), "Should see home content [empty] in body response.")
		}
	}
}

func TestHandlePost(t *testing.T) {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test handling post request")
	{
		t.Log("Test 0: When requesting a post page with valid data")
		{

			r := httptest.NewRequest(http.MethodGet, "/post", nil)
			r.SetPathValue("fn", "lorem-Ipsum")
			r.SetPathValue("dir", "2024")

			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handlePost("content/test/posts").ServeHTTP(w, r)
			testkit.Check(t, strings.Contains(w.Body.String(), "14 Apr 2024"), "Should find date \"14 Apr 2024\" in body response.")
		}

		t.Log("Test 1: When requesting a post page with invalid data")
		{

			r := httptest.NewRequest(http.MethodGet, "/post", nil)
			r.SetPathValue("fn", "lorem-Ipsun") // wrong name
			r.SetPathValue("dir", "2024")

			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handlePost("content/test/posts").ServeHTTP(w, r)
			testkit.Check(t, strings.Contains(w.Body.String(), ErrIBrokeSomething.Error()), "Should find date \"14 Apr 2024\" in body response.")
		}
	}
}

func TestHandleAbout(t *testing.T) {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test handling about request")
	{
		t.Log("\t0:\tWhen requesting about page with valid data")
		{
			want := "about lorem ipsum"

			r := httptest.NewRequest(http.MethodGet, "/about", nil)
			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handleAbout("content/test").ServeHTTP(w, r)

			testkit.Check(t, strings.Contains(w.Body.String(), want), "Should find %s in body response.", want)
		}

		t.Log("\t1:\tWhen requesting about page with valid data")
		{
			want := ErrIBrokeSomething.Error()
			r := httptest.NewRequest(http.MethodGet, "/about", nil)
			w := httptest.NewRecorder()

			handlers := handlers{templateCache: templateCache}
			handlers.handleAbout("content/test/empty").ServeHTTP(w, r)

			testkit.Check(t, strings.Contains(w.Body.String(), want), "Should find \"%s\" in body response.", want)
		}
	}
}
