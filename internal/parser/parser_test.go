package parser

import (
	blog "github.com/thrcd/thr-blog"
	"github.com/thrcd/thr-blog/internal/testkit"
	"strings"
	"testing"
	"time"
)

func TestParseMarkdown(t *testing.T) {
	t.Log("Testing Markdown parsing.")
	{
		fileWithMeta, err := blog.Cfs.ReadFile("content/test/posts/2024/lorem-Ipsum.md")
		if err != nil {
			testkit.ErrorT(t, "Should read file with metadata. Got Error: %v", err)
		}

		fileWithoutMeta, err := blog.Cfs.ReadFile("content/test/posts/2024/lorem-Ipsum.md")
		if err != nil {
			testkit.ErrorT(t, "Should read file with metadata. Got Error: %v", err)
		}

		want := "<h3>lorem ipsum</h3>"

		t.Logf("Test 0:\tWhen parsing markdown with metadata")
		{
			htmlFromFileWithMeta, err := ParseMarkdown(fileWithMeta)
			if err != nil {
				testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
			}

			received := strings.Trim(htmlFromFileWithMeta.Body, "\n")
			testkit.Check(t, strings.Contains(received, want), "Should contain %s.", want)
		}

		t.Logf("Test 1:\tWhen parsing markdown without metada")
		{
			htmlFromFileWithoutMeta, err := ParseMarkdown(fileWithoutMeta)
			if err != nil {
				testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
			}

			received := strings.Trim(htmlFromFileWithoutMeta.Body, "\n")
			testkit.Check(t, strings.Contains(received, want), "Should return html %s.", want)
		}

		t.Logf("Test 2:\tWhen parsing empty markdown")
		{
			fileEmpty := []byte(``)

			htmlEmpty, err := ParseMarkdown(fileEmpty)
			if err != nil {
				testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
			}

			received := strings.Trim(htmlEmpty.Body, "\n")
			successArgs := []any{"", received}
			testkit.Check(t, received == "", "Should return html %s. Received: %s", successArgs...)
		}
	}
}

func TestParseMetadata(t *testing.T) {
	t.Log("Testing markdown parsing and getting metadata.")
	{
		fileWithMeta, err := blog.Cfs.ReadFile("content/test/posts/2024/lorem-Ipsum.md")
		if err != nil {
			testkit.ErrorT(t, "Should read file with metadata. Got Error: %v", err)
		}

		{
			metadata := ParseMetadata(fileWithMeta)
			wantTitle := "lorem ipsum"
			wantDate, _ := time.Parse("2006-01-02", "2024-04-14")
			wantTags := []string{"tag1", "tag2"}

			titleArgs := []any{wantTitle, metadata.Title}
			testkit.Check(t, metadata.Title == wantTitle, "Should return title %s. Received: %s", titleArgs...)

			dateArgs := []any{wantDate, metadata.Date}
			testkit.Check(t, metadata.Date.Day() == wantDate.Day() && metadata.Date.Month() == wantDate.Month() && metadata.Date.Year() == wantDate.Year(), "Should return date %v. Received: %v", dateArgs...)

			tagsArgs := []any{wantTags, metadata.Tags}
			testkit.Check(t, metadata.Tags[0] == wantTags[0] && metadata.Tags[1] == wantTags[1], "Should return tags %s. Received: %s", tagsArgs...)
		}
	}
}

const (
	mdWithMeta = `+++
title = "markdown"
date = "2024-01-28"
tags = ["tag1","tag2"]
+++

# Hello parser
`

	mdWithoutMeta = `# Hello parser`
)

func BenchmarkParseMetadata(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseMetadata([]byte(mdWithMeta))
	}
}

func BenchmarkParseMarkdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseMarkdown([]byte(mdWithMeta))
	}
}

func BenchmarkParseMarkdownWithoutMeta(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseMarkdown([]byte(mdWithoutMeta))
	}
}
