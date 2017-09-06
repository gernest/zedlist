// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package tmpl implements template rendering, from embedded assets.
package tmpl

import (
	"bytes"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/labstack/echo"

	"github.com/Unknwon/com"
	asset "github.com/gernest/zedlist/bindata/template"
	"github.com/gernest/zedlist/modules/i18n"
	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/mdown"
)

const (
	//LoginTpl renders the login page
	LoginTpl = "auth/login.html"

	//RegisterTpl renders registration page
	RegisterTpl = "auth/register.html"

	DeleteTpl = "auth/delete.html"

	JobsNewTpl  = "jobs/new.html"
	JobsViewTpl = "jobs/view.html"

	//RegisterScripts addition scripts to be included in the registration page
	// this adds link to the pickadate.js jquery and other date related javascript.
	RegisterScripts = "/auth/date_picker.html"

	//DashHomeTpl renders dashboard home page
	DashHomeTpl = "dash/home.html"

	//DashJobTpl renders dashboard  jobs home.
	DashJobTpl = "dash/jobs_new.html"

	//DashProfileTpl renders user's profile page
	DashProfileTpl = "dash/profile.html"

	//ResumeHomeTpl renders resume home page
	ResumeHomeTpl = "resume/home.html"

	//ResumeUpddateTpl renders resume update page
	ResumeUpddateTpl = "resume/update.html"

	//ResumeViewTpl renders resume view template.
	ResumeViewTpl = "resume/view.html"

	//BaseHomeTpl renders zedlist home page
	BaseHomeTpl = "base/home.html"

	//BaseJobsHomeTpl renders zedlist's jobs home page
	BaseJobsHomeTpl = "base/jobs.html"

	//BaseJobsViewTpl renders a sinlge job view for zedlist
	BaseJobsViewTpl = "base/jobs_view.html"

	//BaseDocsHomeTpl renders the home page for zedlist docs
	BaseDocsHomeTpl = "base/docs_index.html"

	//BaseDocTpl renders a single zedlist document
	BaseDocTpl = "base/docs.html"

	//BaseRegionsTpl renders zedlist regions home page
	BaseRegionsTpl = "base/regions.html"

	//BaseRegionsJobViewTpl renders jobs by a specified region
	BaseRegionsJobViewTpl = "base/regions_job.html"

	//BaseRegionsPaginateTpl renders jobs by region with pagination
	BaseRegionsPaginateTpl = "base/regions_paginate.html"

	//ErrNotFoundTpl renders 404
	ErrNotFoundTpl = "errors/404.html"

	//ErrServerTpl renders 500
	ErrServerTpl = "errors/500.html"

	//ErrBadRequest renders 400
	ErrBadRequest = "errors/400.html"
)

var (

	// NotFoundMessage is the default 404 message
	NotFoundMessage = "we cant find the resource you asked for, please try again with different details"

	//ServerErrorMessage is the default 500 message
	ServerErrorMessage = "Oops, there was a problem"

	//BadRequestMessage is the default 400 message
	BadRequestMessage = "bad request"
)

// TPL  is the default template object, with all templates loaded
var TPL *Template

func init() {
	config := &Config{
		Name: "base",
		IncludesDirs: []string{
			"base",
			"partials",
			"auth",
			"dash",
			"errors",
			"resume",
			"jobs",
		},
	}
	t, err := New(config)
	if err != nil {
		log.Error(nil, err)
	}
	TPL = t
}

// Funcs is a map of default template functions.
var Funcs = template.FuncMap{
	"date":       date,
	"tr":         translate,
	"md":         toMarkdown,
	"plain":      toHTML,
	"label":      label,
	"script":     script,
	"snake":      com.ToSnakeCase,
	"dashed":     dashed,
	"flag":       flagClass,
	"switchLang": switchLang,
	"username":   useername,
}

// Config is the template configuration. Templates are loaded from embedded source, this act as a
// guide on which to include in the Template.
type Config struct {

	// Name is the name of the base template. All other templates
	// will be parsed and associated with this template.
	Name string

	// Directories to include template files. These should be the directories names
	// under the the templates path. They are relatiove to the template path.
	//
	// Say you want to include the templates in the path /templates/base and /templates/auth.
	// You will need to omit the prefix /templates/ and do like this.
	//		[]string{"base","auth"}
	IncludesDirs []string
}

// Template contains templates that are loaded from embedded assets
type Template struct {
	cfg *Config
	tpl *template.Template
}

// New returns a new template
func New(cfg *Config) (*Template, error) {
	t := &Template{cfg: cfg}
	return t.load()
}

// Render renders a template with name tpl, passing val as context data.
func (t *Template) Render(tpl string, val interface{}) (string, error) {
	out := &bytes.Buffer{}
	err := t.tpl.ExecuteTemplate(out, tpl, val)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// RenderTo renders the templaet with name name and passes data as contx, the result is written
// to out.
func (t *Template) RenderTo(out io.Writer, name string, data interface{}) error {
	err := t.tpl.ExecuteTemplate(out, name, data)
	if err != nil {
		log.Error(nil, err)
		return err
	}
	return nil
	//return t.tpl.ExecuteTemplate(out, name, data)
}

// loads the templates, using constraints specified in the cfg provided at initialization.
// The templates are read from the embedded assets, no extesion observation is made. The
// global Funcs is used as the template funcmap.
func (t *Template) load() (*Template, error) {
	l := t.cfg.Name
	base := template.New(l).Funcs(Funcs)
	t.tpl = base
	var lerr error
	for _, dir := range t.cfg.IncludesDirs {
		for _, name := range asset.AssetNames() {
			if strings.HasPrefix(name, dir) {
				tpl := t.tpl.New(name)
				nd, err := asset.Asset(name)
				if err != nil {
					lerr = err
					break
				}
				_, err = tpl.Parse(string(nd))
				if err != nil {
					lerr = err
					break
				}
			}
		}
	}
	if lerr != nil {
		return nil, lerr
	}
	return t, nil
}

// default date format for zedlist.
func date(t time.Time) string {
	return t.Format(time.RFC822)
}

func flagClass(lang string) template.CSS {
	switch lang {
	case "en":
		return template.CSS("us")
	case "sw":
		return template.CSS("tz")
	default:
		return template.CSS("")
	}
}

func switchLang(lang string) template.CSS {
	switch lang {
	case "en":
		return template.CSS("sw")
	case "sw":
		return template.CSS("en")
	default:
		return template.CSS("en")
	}
}

// translate the str, into a given lang. Only two languages are supported english( en )
// and swahili( sw )
func translate(lang, str string, a ...interface{}) string {
	l := i18n.CloneLang()
	switch lang {
	case "sw":
		l.SetTarget("sw")
		return l.T(str, a...)
	case "en":
		l.SetTarget("en")
		return l.T(str, a...)
	}
	l.SetTarget("en")
	return l.T(str, a...)
}

// renders a given src to github flavored markdown.
func toMarkdown(src string) template.HTML {
	return template.HTML(mdown.Markdown([]byte(src)))
}

// outputs raw html for a given string( It wont be escaped ).
func toHTML(s string) template.HTML {
	return template.HTML(s)
}

// label reformat the rendering of form labels. By default gforms.ModelInstance fields
// are snake cased, so we remove all undescores and put spaces instead. This doesn't care
// about the number of occurance of the underscore.
//
// Example
//		label("mwanza_tanzania") //=> "mwanza tanzania"
func label(s string) string {
	return strings.Replace(s, "_", " ", -1)
}

func useername(s string) string {
	return "@" + s
}

// adds the content of the given file into the template document. This is a helper for
// adding scripts to the templates. The name should be a filepath to the template(that
// whose content we want to insert ).
//
// Example calling this func with name set to /desk/auth/date_picker.html will result
// in the file desk/auth/date_picker.html being injected into the caller template, the name
// is relative to the template root directory which in our case is templates.
func script(name string) template.HTML {
	n := strings.TrimPrefix(name, "/")
	b, err := asset.Asset(n)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(b)
}

// returns a string, which is the result of replacing any underscore character with a dash.
func dashed(str string) string {
	return strings.Replace(str, "_", "-", -1)
}

//Renderer implements echo.Renderer, it renders zedlist templates.
type Renderer struct {
	tp *Template
}

//NewRenderer creates a new renderer instance
func NewRenderer() *Renderer {
	return &Renderer{tp: TPL}
}

//Render renders a template.
func (r *Renderer) Render(out io.Writer, name string, data interface{}, ctx echo.Context) error {
	return r.tp.RenderTo(out, name, data)
}
