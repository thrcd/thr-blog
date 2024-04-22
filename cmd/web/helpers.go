package main

import (
	blog "github.com/thrcd/thr-blog"
	"github.com/thrcd/thr-blog/internal/parser"
	"sort"
	"strings"
)

func getSubDirs(root string) []string {
	dirs, err := blog.Cfs.ReadDir(root)
	if err != nil {
		return nil
	}

	postsDirs := make([]string, 0)
	for _, dir := range dirs {
		if dir.IsDir() {
			url := root + "/" + dir.Name()
			postsDirs = append(postsDirs, url)
		}
	}

	return postsDirs
}

func getFilePaths(dir string) ([]string, error) {
	files, err := blog.Cfs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, dir+"/"+file.Name())
		}
	}

	return paths, nil
}

func lastSubString(path, delim string) string {
	i := strings.LastIndex(path, delim)
	return path[i+1:]
}

func maps[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func getPostListItems(dirs []string) map[string][]postListItem {
	collection := make(map[string][]postListItem)

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] > dirs[j]
	})

	for _, dir := range dirs {
		paths, err := getFilePaths(dir)
		if err != nil {
			return nil
		}
		curDir := lastSubString(dir, "/")

		posts := make([]postListItem, 0)
		for _, path := range paths {
			file, err := blog.Cfs.ReadFile(path)
			if err != nil {
				return nil
			}

			filename := lastSubString(path, "/")
			metadata := parser.ParseMetadata(file)

			postItem := postListItem{
				Filename: filename,
				Metadata: metadata,
			}

			posts = append(posts, postItem)
		}

		sort.Slice(posts, func(i, j int) bool {
			return posts[j].Metadata.Date.Before(posts[i].Metadata.Date)
		})

		collection[curDir] = posts
	}

	return collection
}

func getMarkdown(filePath string) (parser.Markdown, error) {
	file, err := blog.Cfs.ReadFile(filePath)
	if err != nil {
		return parser.Markdown{}, err
	}

	md, err := parser.ParseMarkdown(file)
	if err != nil {
		return parser.Markdown{}, err
	}

	return md, nil
}
