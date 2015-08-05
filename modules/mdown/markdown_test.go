package mdown

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func ExampleCreate() {
	anchorName := CreateAnchor("This is a header")

	fmt.Println(anchorName)

	// Output:
	// this-is-a-header
}

func ExampleCreate2() {
	fmt.Println(CreateAnchor("This is a header"))
	fmt.Println(CreateAnchor("This is also          a header"))
	fmt.Println(CreateAnchor("main.go"))
	fmt.Println(CreateAnchor("Article 123"))
	fmt.Println(CreateAnchor("<- Let's try this, shall we?"))
	fmt.Printf("%q\n", CreateAnchor("        "))
	fmt.Println(CreateAnchor("Hello, 世界"))

	// Output:
	// this-is-a-header
	// this-is-also-a-header
	// main-go
	// article-123
	// let-s-try-this-shall-we
	// ""
	// hello-世界
}

func ExampleMarkdown() {
	text := []byte("Hello world github/linguist#1 **cool**, and #1!")

	os.Stdout.Write(Markdown(text))

	// Output:
	//<p>Hello world github/linguist#1 <strong>cool</strong>, and #1!</p>
}

// An example of how to generate a complete HTML page, including CSS styles.
func ExampleMarkdown_completeHtmlPage() {
	var w io.Writer = os.Stdout // It can be an http.ResponseWriter.
	markdown := []byte("# GitHub Flavored Markdown\n\nHello.")

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href=".../github-flavored-markdown.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)

	// Output:
	//<html><head><meta charset="utf-8"><link href=".../github-flavored-markdown.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;"><h1><a name="github-flavored-markdown" class="anchor" href="#github-flavored-markdown" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>GitHub Flavored Markdown</h1>
	//
	//<p>Hello.</p>
	//</article></body></html>
}

func ExampleHeader() {
	text := []byte("## git diff")

	os.Stdout.Write(Markdown(text))

	// Output:
	//<h2><a name="git-diff" class="anchor" href="#git-diff" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>git diff</h2>
}

func TestBlockQuote(t *testing.T) {
	in, err := ioutil.ReadFile("test/input.md")
	if err != nil {
		t.Error(err)
	}
	out := Markdown(in)
	expect, err := ioutil.ReadFile("test/out.html")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(out, expect) {
		t.Errorf(" expected %s got %s", expect, out)
	}
}
