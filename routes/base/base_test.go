// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/Unknwon/com"
	userAuth "github.com/gernest/zedlist/middlewares/auth"
	"github.com/gernest/zedlist/middlewares/flash"
	"github.com/gernest/zedlist/middlewares/i18n"
	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
	"github.com/labstack/echo"
)

var ts = testServer()
var client = &http.Client{Jar: getJar()}

func testServer() *httptest.Server {
	e := echo.New()
	e.Renderer = tmpl.NewRenderer()

	// middlewares
	e.Use(utils.WrapMiddleware(i18n.Langs()))      // languages
	e.Use(utils.WrapMiddleware(flash.Flash()))     // flash messages
	e.Use(utils.WrapMiddleware(userAuth.Normal())) // adding user context data

	// HOME
	e.GET("/", Home)
	e.GET("/language/:lang", SetLanguage)
	b := e.Group("/jobs")
	b.GET("/", JobsHome)
	b.GET("/view/:id", JobView)
	b.GET("/regions", RegionsHome)
	b.GET("/regions/:name", RegionsJobView)
	b.GET("/regions/:name/:from/:to", RegionsJobPaginate)

	// DOCS
	e.GET("/docs", DocsHome)
	e.GET("/docs/:name", Docs)

	return httptest.NewServer(e)
}

func getJar() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Log
	}
	return jar
}

func closeTest() {
	ts.Close()
}

//
//
//	BASE ROUTES
//
//
var JobsSample = []struct {
	title, desc string
}{
	{"Home made robots expert", "roboting"},
	{"financier", " a fancy finance maestro"},
	{"dreadful fan", "where the heck is sembene?"},
	{"bongo fleva fan", "suckers"},
}

var regionsSample = []struct {
	name, short string
}{
	{"mwanza", "mza"},
	{"dar es salaam", "dar"},
}

func TestHome(t *testing.T) {

	// Migrate, this is going to be used by all subsequent tests.
	for _, v := range regionsSample {
		q := db.Conn.FirstOrCreate(&models.Region{}, models.Region{Name: v.name, Short: v.short})
		if q.Error != nil {
			t.Errorf("migrating regios %v", q.Error)
		}
	}
	regs, err := query.GetAllRegions(db.Conn)
	if err != nil {
		t.Errorf("retrieving regions %v", err)
	}
	for _, v := range regs {
		for _, job := range JobsSample {
			j := &models.Job{}
			j.Title = job.title
			j.Description = job.desc
			j.Region = *v
			q := db.Conn.Create(j)
			if q.Error != nil {
				t.Errorf("creating jobs %v", q.Error)
			}
		}
	}

	b, err := com.HttpGetBytes(client, ts.URL, nil)
	if err != nil {
		t.Errorf("getting home page %v", err)
	}

	// Yes a home page should contain zedlist tittle
	if !bytes.Contains(b, []byte("<title>zedlist</title>")) {
		t.Errorf("expected home page got %s", b)
	}
}

func TestSetLang(t *testing.T) {
	sw := fmt.Sprintf("%s/language/sw", ts.URL)
	b, err := com.HttpGetBytes(client, sw, nil)
	if err != nil {
		t.Errorf("getting home page %v", err)
	}

	// Yes a home page should contain swahili words
	if !bytes.Contains(b, []byte("Nyumbani")) {
		t.Errorf("expected home page to be in swahili got %s", b)
	}
}

func TestJobsHome(t *testing.T) {
	home := fmt.Sprintf("%s/jobs/", ts.URL)
	b, err := com.HttpGetBytes(client, home, nil)
	if err != nil {
		t.Errorf("getting home page %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		t.Errorf("loading document %v", err)
	}
	s := doc.Find(".job-item")
	if s.Length() < 8 {
		t.Errorf("expected 8 jobs got %d", s.Length())
	}
}

func TestJobsView(t *testing.T) {
	view := "%s/jobs/view/%d"
	jobs, err := query.GetLatestJobs(db.Conn)
	if err != nil {
		t.Error(err)
	}
	for _, job := range jobs {
		viewURL := fmt.Sprintf(view, ts.URL, job.ID)
		resp, err := client.Get(viewURL)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, resp.StatusCode)
		}
	}
}
func TestJobsRegionsHome(t *testing.T) {
	home := fmt.Sprintf("%s/jobs/regions", ts.URL)
	b, err := com.HttpGetBytes(client, home, nil)
	if err != nil {
		t.Errorf("getting regions home page %v", err)
	}
	if !bytes.Contains(b, []byte(regionsSample[1].name)) {
		t.Errorf("expected %s to contain %s", b, regionsSample[1].name)
	}
}

func TestJobsRegionsByShortName(t *testing.T) {
	regs, err := query.GetAllRegions(db.Conn)
	if err != nil {
		t.Errorf("retriving regions %v", err)
	}
	for _, v := range regs {
		regHome := fmt.Sprintf("%s/jobs/regions/%s", ts.URL, v.Short)
		b, err := com.HttpGetBytes(client, regHome, nil)
		if err != nil {
			t.Errorf("getting regions home page %v", err)
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
		if err != nil {
			t.Errorf("loading document %v", err)
		}
		s := doc.Find(".job-item")
		if s.Length() != 4 {
			t.Errorf("expected 4 got %d", s.Length())
		}
	}
}

func TestJobsRegionsPaginate(t *testing.T) {
	pg := []struct {
		from, to int
	}{
		{0, 1},
		{0, 2},
		{1, 4},
	}
	for _, page := range pg {
		for _, reg := range regionsSample {
			paginate := fmt.Sprintf("%s/jobs/regions/%s/%d/%d", ts.URL, reg.short, page.from, page.to)
			b, err := com.HttpGetBytes(client, paginate, nil)
			if err != nil {
				t.Errorf("getting regions home page %v", err)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
			if err != nil {
				t.Errorf("loading document %v", err)
			}
			s := doc.Find(".job-item")
			d := page.to - page.from
			if s.Length() != d {
				t.Errorf("expected %d got %d", d, s.Length())
			}
		}

	}
}

//
//
//		DOCS
//
//
func TestDocsHome(t *testing.T) {
	l := fmt.Sprintf("%s/docs", ts.URL)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting docs home %v", err)
	}
	if !bytes.Contains(b, []byte("home.md")) {
		t.Errorf("ecpectd docs home got %s", b)
	}
}
func TestDocs(t *testing.T) {

	// with .md  extension
	l := fmt.Sprintf("%s/docs/home.md", ts.URL)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting docs home %v", err)
	}
	if !bytes.Contains(b, []byte("home.md")) {
		t.Errorf("ecpectd docs home got %s", b)
	}

	// without .md extsnison
	l = fmt.Sprintf("%s/docs/home", ts.URL)
	b, err = com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting docs home %v", err)
	}
	if !bytes.Contains(b, []byte("home.md")) {
		t.Errorf("ecpectd docs home got %s", b)
	}
}

// this should be always at the bottom of this file.
func TestCLose(t *testing.T) {
	closeTest()
}
