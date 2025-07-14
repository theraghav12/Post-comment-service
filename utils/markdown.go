package utils

import (
	bytes "bytes"
	"github.com/yuin/goldmark"
)

// RenderMarkdown converts markdown string to HTML string
func RenderMarkdown(input string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(input), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
