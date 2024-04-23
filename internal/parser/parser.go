package parser

import (
	"bufio"
	"bytes"
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

func ParseMarkdown(file []byte) (Markdown, error) {
	var markdown Markdown
	var buf bytes.Buffer

	// Checks if file has metadata. If positive, it parses the metadata and assign to markdown.
	if len(file) != 0 && metaDelimReg.Match(file[:1]) {
		//var inMetadata bool
		var line []byte
		var lastDelimIndex int
		var metaDelimCounter int

		markdown.Metadata = ParseMetadata(file)

		reader := bytes.NewReader(file)
		scanner := bufio.NewScanner(reader)

		// Scans the files and increments the lastDelimIndex counter until it finds the last occurrence of "+++".
		// This counter will be used to skip over the metadata when parsing the markdown.
		for i := 0; scanner.Scan(); i++ {
			line = bytes.TrimSpace(scanner.Bytes())
			lastDelimIndex += len(line) + len("\n")

			if metaDelimReg.Match(line) {
				metaDelimCounter++
			}

			if metaDelimCounter > 1 {
				break
			}
		}

		// File now has only the markdown content without metadata
		file = file[lastDelimIndex:]

		if err := scanner.Err(); err != nil { // Checked scanner error
			return Markdown{}, err
		}
	}

	var bodyBuf bytes.Buffer
	if buf.Len() == 0 {
		buf.Write(file)
	}

	md := blackfriday.Run(file)
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
