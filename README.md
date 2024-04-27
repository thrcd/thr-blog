
# Thr Blog

> Simple personal blog written in Go.

_the readme is almost a blog post_

I've been writing my whole life. Most of the time, I write on paper, but I also have a blog-like setup running on my Raspberry Pi for the last 5-6 years. This is a place where I store everything that's on my mind to read later.

Over the past month, I've rediscovered the habit I lost in 2010: reading blogs. Specifically, simple blogs without the nonsense of Medium, pop-ups, tracking, ADS... just words and personality.

So, one Saturday afternoon, a typical cold winter day in Wisconsin, after giving up for the third time on reading "Crime and Punishment", I decided to write a simple blog website using Go, where I could just dump my .md files into a folder.

Why not Hugo? I didn't know about Hugo, but it was all I could find after a quick search for "how to parse markdown in Go" and "blog in go". After a quick look at their repo, I think it's great, but I wanted to build something for fun, something simple I can explain and understand from the beginning. Why use a sledgehammer to crack a nut, right?

## Content

The content of the blog is structured by folders. Right now, the blog only has two pages: posts and about.

Every folder inside the "posts" folder will be a section on the posts page. And every .md file inside the subfolders will be a post.

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

Then, I start every Markdown file with a simple metadata section:

```
+++
title = "lorem ipsum"
date = "2024-04-14" 
tags = ["some tag", "other tag"]
+++

... Post

```


## Run Locally

Clone the project

```bash
  git clone https://github.com/thrcd/thr-blog.git
```

Go to the project directory

```bash
  cd thr-blog
```

There is a simple makefile to help run and build. 

Run tests

```bash
  make tests
```

or 

```bash
  go test ./...
```

Build

```bash
  make build
```
or 

Build

```bash
  go build ./cmd/web
```

Run (dev)

```bash
  make dev
```

or 

```bash
  go run ./cmd/web
```

## Acknowledgements

At the beginning, I was writing my own Markdown parser just for fun, but with a new baby at home, time has become more unpredictable, so it's still a work in progress—maybe forever. I try to use third-party libraries as little as possible, or at least until I understand how they work and how they might break in the future. After another afternoon reading the code from blackfriday, I started using this amazing Markdown parser from russross to parse the post section of my .md files.

 - [blackdriday](github.com/russross/blackfriday/v2)
## Todo

- [ ]  Actually write and post something
- [x]  Brew some coffee
- [x]  Relax  
