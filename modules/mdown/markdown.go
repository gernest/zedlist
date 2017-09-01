package mdown

import "github.com/gernest/github_flavored_markdown"

// Markdown renders GitHub Flavored Markdown text.
func Markdown(text []byte) []byte {
	return github_flavored_markdown.Markdown(text)
}
