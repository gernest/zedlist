/*
Package mdown contains functions for rendering Github flavored narkdown. Majority of this
code is borowwed from https://github.com/shurcooL/go/github_flavored_markdown. Some aspects have
been changed to meet my needs.
*/
package mdown

import "github.com/gernest/github_flavored_markdown"

// Markdown renders GitHub Flavored Markdown text.
func Markdown(text []byte) []byte {
	return github_flavored_markdown.Markdown(text)
}
