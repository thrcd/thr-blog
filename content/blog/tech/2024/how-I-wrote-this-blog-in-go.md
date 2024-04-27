+++
title = "How I wrote this blog in go"
date = "2024-04-22"
tags = ["Go version=1.22.2"]
+++

I've been writing my whole life. Most of the time, I write on paper, but I also have a blog-like setup running on my Raspberry Pi for the last 5-6 years. This is a place where I store everything that's on my mind to read later.

Over the past month, I've rediscovered the habit I lost in 2010: reading blogs. Specifically, simple blogs without the nonsense of Medium, pop-ups, tracking, ADS... just words and personality.

So, one Saturday afternoon, a typical cold winter day in Wisconsin, after giving up for the third time on reading "Crime and Punishment", I decided to write a simple blog website using Go, where I could just dump my .md files into a folder.

Why not Hugo? I didn't know about Hugo, but it was all I could find after a quick search for "how to parse markdown in Go" and "blog in go". After a quick look at their repo, I think it's great, but I wanted to build something for fun, something simple I can explain and understand from the beginning. Why use a sledgehammer to crack a nut, right?

I basically wanted to write something right on my text editor, save as .md inside one folder, type-enter git commit-push, and the post was online. The idea was to have a content folder > a posts folder > year folder > post files .md.

The final structure ended up like this:

```
* content folder:
  - blog:
      about.md
      home.md
      - tech:
          - 2024:
              - tech post files *.md
      - life:
          - 2024:
              - life post files *.md
  - tests (files to run the tests mocking the blog files)
```

When it comes to project structure in Go, there's no one-size-fits-all guideline like in some other languages.

For my **(small)** projects, I often keep all the files in a single folder. However, as the project grows, I might organize them into packages. The structure usually reflects the project's needs, my design choices for decoupling, and sometimes even my mood.

Still, there are some common folders I tend to use. cmd, internal, pkgs, and so on.

```
* blog:
    - cmd
        - web: main files (main.go, handlers, helpers, template ...)
    - internal: core pkgs (parser, ui, testkit)
    - content: posts and page content
```

The blog is divided into four main sections:

```
(This was supposed to be a pyramid.)

         -- server -- 
               ↑
    --- handlers, routes ---
               ↑
  --------- helpers ----------
               ↑
---------- parser, ui -----------
```

### parser and ui

I'll start from the bottom up. The ui folder houses all the Go template files (.tmpl), CSS, and images. The [embed](https://pkg.go.dev/embed) package simplifies handling these files.

```go
// mfs.go
// This is located at ui package and embeds all html and static files

//go:embed "html" "static"
var Files embed.FS
var StaticFS, _ = fs.Sub(Files, "static")

// cfs.go
// Go file at root folder and embeds the content folder, holding all the posts.
// I'll probably move this to an external storage later. ~maybe~.

//go:embed "content"
var Cfs embed.FS
```

Then I have the parser package. At the beginning, I was writing my own Markdown parser just for fun, but with a new baby at home, time has become more unpredictable, so it's still a work in progress—maybe forever.

I try to use third-party libraries as little as possible, or at least until I understand how they work and how they might break in the future. After another afternoon reading the code from [blackdriday](https://github.com/russross/blackfriday), I started using this amazing Markdown parser from russross to parse the post section of my .md files.

I still have a parser pkg to do custom stuff. For example, I use the metadata in the .md file in TOML format, so I wrote my own custom code to handle it.

First, I defined the custom type to represent the markdown and the metadata.

```go
// parser.go

type Metadata struct {
    Title string
    Date  time.Time
    Tags  []interface{}
}

type Markdown struct {
    Metadata Metadata
    Body     string
}
```

Then I expose two functions:

```
// parser.go

// Custom code to parse only the metadata
// Basically, ParseMetada scans the file ([]bytes) 
// and uses regexp.MustCompile to get the metadata from
// file, and breakers the scan when it finds the second delimeter (+++).  
func ParseMetadata(b []byte) Metadata 

// Uses blackdriday lib to parse the content of the markdown.
func ParseMarkdown(b []byte) (Markdown, error) 
```

_I won't paste the implementation here; otherwise, this post would become much longer than I'd like. But you can check them out in the repo._

### helpers

There's no need to overcomplicate things; remember, less is more, so let's keep it simple. I have a file called helpers.go that contains a few helper functions to construct the data used by the templates. Let's go through those.

```go
// helpers.go

// getSubDirs returns a slice of subdirectories inside a given directory.
// For example, if root := "content/posts", 
// it will return ["content/posts/2024", "content/posts/2023", ...].
func getSubDirs(root string) []string 

// getFilePaths returns a slice of file paths within a specific directory.
// For example, if dir := "content/posts/2024", 
// it returns ["content/posts/2024/lorem-Ipsum.md"].
func getFilePaths(dir string) ([]string, error)

// lastSubString returns the last substring in a given path
// delimited by a specified delimiter.
// For example, if path := "content/posts/2024/lorem-Ipsum.md" 
// and delim := "/", it returns "lorem-Ipsum.md".
func lastSubString(path, delim string) string

// maps is a simple map function implementation.
func maps[T, V any](ts []T, fn func(T) V) []V

// getPostListItems constructs data for the post list page.
// Each postListItem contains Filename and Metadata.
func getPostListItems(dirs []string) map[string][]postListItem 

// getMarkdown reads a file, calls the parser, and returns the markdown content.
// This function is used in post and about handlers.
func getMarkdown(filePath string) (parser.Markdown, error)
```

These functions manage the entire data flow from the `.md` files to the Go templates.

### handlers

I've got three handlers: one for handling posts on the homepage, another for individual posts, and one for the about page.

I've also defined a "handlers" type to pass along any shared dependencies that my handlers might require. At the moment, the only dependency it receives is templateCache, which is a map to cache all the Go templates.

```go
// handlers.go

type handlers struct {
    templateCache map[string]*template.Template
}
```

Next, I have my first handler called handlePosts. I prefer using this naming convention and returning an http.HandlerFunc rather than the typical postHandler(w http.ResponseWriter, r *http.Request).

This approach makes it easier for me to read the code and pass additional dependencies, and I've found it simplifies testing these handlers as well.

In this handler, I return an anonymous function that implements (w http.ResponseWriter, r *http.Request). The rest of the code involves calling several helper functions: one to fetch the list of directories inside "content > post", another to get the list of posts, and one to construct the template data.

This data function stores the necessary information for my templates. Finally, I render the post.tmpl template.

```go
// handlers.go

func (h *handlers) handlePosts(root string) http.HandlerFunc
    return func(w http.ResponseWriter, r *http.Request) {
        postType := strings.TrimPrefix(r.RequestURI, "/")
        
        dirs := getSubDirs(root + "/" + postType)
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
        data.PostType = postType
        
        render(w, "posts.tmpl", h.templateCache, data)
    }
}
```

Oh... There are other parts I won't delve into here, but you can explore them in the repository. The "empty" handler method renders `empty.tmpl` if there are no posts in the posts folder.

Also there is the "render", which fetches the template and executes it.

*You can find additional handlers and more details in the repository.*

One more thing to note in this section is the routing. It utilizes the Go standard library and is quite straightforward. It initializes the cache for templates, sets up the handlers struct, and defines the routes.

```go
// routes.go

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
```

### server

Finally, on main.go I create a straightforward server that does exactly what I need. It begins listening and serving.

The server is given the TCP address to listen on, which here is ":4000", and a routes function that returns an http.Handler.

```go
// main.go

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
```

In the future ~probably never~, if I want to add another post section, such as 'coffee' or 'music,'
all I have to do is create a new folder in content/blog/, like content/blog/music/2024/, 
and add a new route to the routes function. There are hundreds of smarter ways to do this, 
but this was the best I could do with the time available. It's working, it's simple, and I like it.

After all this, you run the command:

```bash
go run ./cmd/web
```

and voilà, the blog is alive.

```json
{"time":"2024-04-22T13:42:17.306832-05:00","level":"INFO","msg":"starting server","addr":":4000"}
```

Now the hard part: *writing*.

github repo: [thr-blog](https://github.com/thrcd/thr-blog)