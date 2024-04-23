package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"regexp"
	"strings"
	"time"
)

var (
	metaTitleReg = regexp.MustCompile(`(?i)title\s*=\s*"([^"]+)"`)
	metaDateReg  = regexp.MustCompile(`(?i)date\s*=\s*"([^"]+)"`)
	metaTagsReg  = regexp.MustCompile(`tags\s*=\s*\[([^\]]*)\]`)
	metaDelimReg = regexp.MustCompile(`\++`)
)

type Metadata struct {
	Title string
	Date  time.Time
	Tags  []interface{}
}

type Markdown struct {
	Metadata Metadata
	Body     string
}

func ParseMetadata(b []byte) Metadata {
	reader := bytes.NewReader(b)
	var metadata Metadata
	metaDelimCounter := 0

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		if metaDelimReg.Match(line) {
			metaDelimCounter++
		}

		if metaDelimCounter > 1 {
			break
		}

		ok, titleMatch := getMetadataString(metaTitleReg, line)
		if ok {
			metadata.Title = titleMatch
		}

		ok, dateStr := getMetadataString(metaDateReg, line)
		if ok {
			date, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				metadata.Date = date
			}
		}

		ok, tagsMatch := getMetadataSlice(metaTagsReg, line)
		if ok {
			metadata.Tags = tagsMatch
		}
	}

	if metadata.Date.IsZero() {
		metadata.Date = time.Now()
	}

	return metadata
}

func ParseMarkdown(b []byte) (Markdown, error) {
	var markdown Markdown
	var buf bytes.Buffer

	// Checks if file has metadata. If positive, it parses the metadata and assign to markdown.
	if len(b) != 0 && metaDelimReg.Match(b[:1]) {
		//var inMetadata bool
		var line []byte
		var lastDelimIndex int

		metaDelimCounter := 0

		markdown.Metadata = ParseMetadata(b)

		reader := bytes.NewReader(b)
		scanner := bufio.NewScanner(reader)

		// Scans the files and increments the lastDelimIndex counter until it finds the last occurrence of "+++".
		// This counter will be used to skip over the metadata when parsing the markdown.
		for i := 0; scanner.Scan(); i++ {
			line = bytes.TrimSpace(scanner.Bytes())

			for i := 0; i < len(line); i++ {
				lastDelimIndex++
				fmt.Println(lastDelimIndex)
			}

			if metaDelimReg.Match(line) {
				metaDelimCounter++
			}

			if metaDelimCounter > 1 {
				break
			}
		}

		b = b[lastDelimIndex+1:]

		if err := scanner.Err(); err != nil { // Checked scanner error
			return Markdown{}, err
		}
	}

	var bodyBuf bytes.Buffer
	if buf.Len() == 0 {
		buf.Write(b)
	}

	md := blackfriday.Run(b)
	bodyBuf.Write(md)
	markdown.Body = bodyBuf.String()

	return markdown, nil
}

func getMetadataString(rgx *regexp.Regexp, b []byte) (bool, string) {
	matches := rgx.FindSubmatch(b)

	if len(matches) > 1 {
		return true, string(matches[1])
	}

	return false, ""
}

func getMetadataSlice(rgx *regexp.Regexp, b []byte) (bool, []interface{}) {
	matches := rgx.FindSubmatch(b)

	sl := make([]interface{}, 0)
	if len(matches) > 1 {
		str := string(matches[1])
		str = strings.ReplaceAll(str, "\"", "")
		strs := strings.Split(str, ",")
		for _, s := range strs {
			sl = append(sl, s)
		}

		return true, sl
	}

	return false, nil
}
