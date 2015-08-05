/*
Package mdown contains functions for rendering Github flavored narkdown. Majority of this
code is borowwed from https://github.com/shurcooL/go/github_flavored_markdown. Some aspects have
been changed to meet my needs.
*/
package mdown

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"

	"github.com/gernest/highlight_go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
)

// CreateAnchor returns a sanitized anchor name for the given text.
func CreateAnchor(text string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range []rune(text) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}

// Markdown renders GitHub Flavored Markdown text.
func Markdown(text []byte) []byte {
	htmlFlags := 0
	renderer := &renderer{Html: blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}

	// Parser extensions for GitHub Flavored Markdown.
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	//extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK

	unsanitized := blackfriday.Markdown(text, renderer, extensions)

	// GitHub Flavored Markdown-like sanitization policy.
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowDataURIImages()

	return p.SanitizeBytes(unsanitized)
}

type renderer struct {
	*blackfriday.Html
}

// Header GitHub Flavored Markdown header with clickable and hidden anchor.
func (r *renderer) Header(out *bytes.Buffer, text func() bool, level int, _ string) {
	marker := out.Len()
	doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}

	textString := out.String()[marker:]
	out.Truncate(marker)

	anchorName := CreateAnchor(html.UnescapeString(textString))

	out.WriteString(fmt.Sprintf(`<h%d><a name="%s" class="anchor" href="#%s" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>`, level, anchorName, anchorName))
	out.WriteString(textString)
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

// TODO: Clean up and improve this code.
// BlockCode GitHub Flavored Markdown fenced code block with highlighting.
func (r *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	doubleSpace(out)

	// parse out the language name
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		out.WriteString(`<div class="highlight highlight-`)
		attrEscape(out, []byte(elt))
		lang = elt
		out.WriteString(`"><pre>`)
		count++
		break
	}

	if count == 0 {
		out.WriteString("<pre><code>")
	}

	if highlightedCode, ok := highlightCode(text, lang); ok {
		out.Write(highlightedCode)
	} else {
		attrEscape(out, text)
	}

	if count == 0 {
		out.WriteString("</code></pre>\n")
	} else {
		out.WriteString("</pre></div>\n")
	}
}

var gfmHTMLConfig = syntaxhighlight.HTMLConfig{
	String:        "s",
	Keyword:       "k",
	Comment:       "c",
	Type:          "n",
	Literal:       "o",
	Punctuation:   "p",
	Plaintext:     "n",
	Tag:           "tag",
	HTMLTag:       "htm",
	HTMLAttrName:  "atn",
	HTMLAttrValue: "atv",
	Decimal:       "m",
}

// TODO: Support highlighting for more languages.
func highlightCode(src []byte, lang string) (highlightedCode []byte, ok bool) {
	switch lang {
	case "go":
		var buf bytes.Buffer
		err := highlight_go.Print(src, &buf, syntaxhighlight.HTMLPrinter(gfmHTMLConfig))
		if err != nil {
			return nil, false
		}
		return buf.Bytes(), true
	case "go-old":
		var buf bytes.Buffer
		err := syntaxhighlight.Print(syntaxhighlight.NewScanner(src), &buf, syntaxhighlight.HTMLPrinter(gfmHTMLConfig))
		if err != nil {
			return nil, false
		}
		return buf.Bytes(), true

	default:
		return nil, false
	}
}

// Unexported blackfriday helpers.

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}
