package parser

import (
	"fmt"
	"github.com/thrcd/thr-blog/internal/testkit"
	"strings"
	"testing"
)

const (
	mdWithMeta = `+++
title = "markdown"
date = "2024-01-28"
tags = ["tag1","tag2"]
+++

# Hello parser
`

	mdWithoutMeta = `# Hello parser`
	mdEmpty       = ``

	expectedHtml = "<h1>Hello parser</h1>"
)

func TestParseMarkdown(t *testing.T) {
	fileWithMeta := []byte(mdWithMeta)
	fileWithoutMeta := []byte(mdWithoutMeta)
	fileEmpty := []byte(mdEmpty)

	t.Log("Testing Markdown parsing.")
	{
		t.Logf("Test 0:\tWhen parsing markdown with metadata")
		htmlFromFileWithMeta, err := ParseMarkdown(fileWithMeta)
		if err != nil {
			testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
		}

		received := strings.Trim(htmlFromFileWithMeta.Body, "\n")
		successArgs := []any{expectedHtml, received}
		testkit.Check(t, received == expectedHtml, "Should return html %s. Received: %s", successArgs...)
	}

	{
		t.Logf("Test 1:\tWhen parsing markdown without metada")
		htmlFromFileWithoutMeta, err := ParseMarkdown(fileWithoutMeta)
		if err != nil {
			testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
		}

		received := strings.Trim(htmlFromFileWithoutMeta.Body, "\n")
		successArgs := []any{expectedHtml, received}
		testkit.Check(t, received == expectedHtml, "Should return html %s. Received: %s", successArgs...)
	}

	{
		t.Logf("Test 2:\tWhen parsing empty markdown")
		htmlEmpty, err := ParseMarkdown(fileEmpty)
		if err != nil {
			testkit.ErrorT(t, "Should parse markdown input. [%s]", err)
		}

		received := strings.Trim(htmlEmpty.Body, "\n")
		successArgs := []any{"", received}
		testkit.Check(t, received == "", "Should return html %s. Received: %s", successArgs...)
	}
}

func TestParseMetadata(t *testing.T) {
	fileWithMeta := []byte(mdWithMeta)

	t.Log("Testing markdown parsing and getting metadata.")
	{
		metadata := ParseMetadata(fileWithMeta)

		titleArgs := []any{"markdown", metadata.Title}
		testkit.Check(t, metadata.Title == "markdown", "Should return title %s. Received: %s", titleArgs...)

		dateArgs := []any{fmt.Sprintf("%d/%d/%d", metadata.Date.Month(), metadata.Date.Day(), metadata.Date.Year())}
		testkit.Check(t, metadata.Date.Day() == 28 && metadata.Date.Month() == 1 && metadata.Date.Year() == 2024, "Should return date 01/28/2024. Received: %s", dateArgs...)

		tagsArgs := []any{"[tag1 tag2]", metadata.Tags}
		testkit.Check(t, metadata.Tags[0] == "tag1" && metadata.Tags[1] == "tag2", "Should return tags %s. Received: %s", tagsArgs...)
	}
}

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
