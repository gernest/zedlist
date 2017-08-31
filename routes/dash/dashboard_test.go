package dash

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Unknwon/com"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	userAuth "github.com/gernest/zedlist/middlewares/auth"
	"github.com/gernest/zedlist/middlewares/flash"
	"github.com/gernest/zedlist/middlewares/i18n"
	"github.com/gernest/zedlist/routes/auth"
	"github.com/labstack/echo"
	"github.com/oxtoacart/bpool"
)

var ts = testServer()
var client = &http.Client{Jar: getJar()}
var dashPath = fmt.Sprintf("%s/dash", ts.URL)
var bufPool = bpool.NewBufferPool(10)

func testServer() *httptest.Server {
	e := echo.New()
	e.Renderer = tmpl.NewRenderer()

	// middlewares
	e.Use(utils.WrapMiddleware(i18n.Langs()))      // languages
	e.Use(utils.WrapMiddleware(flash.Flash()))     // flash messages
	e.Use(utils.WrapMiddleware(userAuth.Normal())) // adding user context data

	// BASE

	// DASHBOARD
	d := e.Group("/dash")
	d.GET("/", Home)
	d.GET("/jobs/new", JobsNewGet)
	d.POST("/jobs/new", JobsNewPost)
	d.GET("/profile", Profile)
	d.POST("/profile/name", ProfileName)

	// AUTH
	e.POST("/login", auth.LoginPost)
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
//		DASHBOARD
//
//
func TestHome(t *testing.T) {
	l := fmt.Sprintf("%s/", dashPath)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting dashboard home %v", err)
	}
	title := "<title>dashboard</title>"
	if !bytes.Contains(b, []byte(title)) {
		t.Errorf(" expected login page got %s", b)
	}
}
func TestJobsNewGet(t *testing.T) {
	l := fmt.Sprintf("%s/jobs/new", dashPath)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting dashboard home %v", err)
	}
	title := "<title>new job</title>"
	if !bytes.Contains(b, []byte(title)) {
		t.Errorf(" expected new job  page got %s", b)
	}
}
func TestJobsNewPost(t *testing.T) {
	loginPath := fmt.Sprintf("%s/login", ts.URL)
	l := fmt.Sprintf("%s/jobs/new", dashPath)

	loginForm := url.Values{
		"username": {"root@home.com"},
		"password": {"superroot"},
	}

	jobForm := url.Values{
		"title":       {"hellow Tanzania"},
		"description": {"my city, my pride"},
	}

	// Case not authorized
	res, err := client.PostForm(l, jobForm)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected %d got %d", http.StatusUnauthorized, res.StatusCode)
	}

	// Authorize the user
	resp, err := client.PostForm(loginPath, loginForm)
	if err != nil {
		t.Errorf("loggging in %v", err)
	}
	defer resp.Body.Close()

	resp1, err := client.PostForm(l, jobForm)
	if err != nil {
		t.Errorf("creating a new job %v", err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp1.StatusCode)
	}

	// case invalid form
	// TODO check whether the redirection was made.
	jobForm.Set("title", "")
	resp2, err := client.PostForm(l, jobForm)
	if err != nil {
		t.Errorf("creating a new job %v", err)
	}
	defer resp2.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp1.StatusCode)
	}

}

func TestProfile(t *testing.T) {
	l := fmt.Sprintf("%s/profile", dashPath)
	resp, err := client.Get(l)
	if err != nil {
		t.Errorf("getting profile page %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.StatusCode)
	}
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	io.Copy(buf, resp.Body)
	title := "<title>profile</title>"
	if !strings.Contains(buf.String(), title) {
		t.Errorf(" expected new job  page got %s", buf.String())
	}
}

func TestProfileName(t *testing.T) {
	l := fmt.Sprintf("%s/profile/name", dashPath)
	nameVars := url.Values{
		"given_name":  {"sharo"},
		"middle_name": {"militon"},
		"family_name": {"aurora"},
	}

	resp, err := client.PostForm(l, nameVars)
	if err != nil {
		t.Errorf("posting form %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.StatusCode)
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)
	io.Copy(buf, resp.Body)

	if !strings.Contains(buf.String(), nameVars.Get("given_ame")) {
		t.Errorf("expected %s to contain %s", buf.String(), nameVars.Get("given_name"))
	}
}

// this should be always at the bottom of this file.
func TestCLose(t *testing.T) {
	closeTest()
}
