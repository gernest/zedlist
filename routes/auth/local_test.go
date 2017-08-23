// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

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

	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/Unknwon/com"
	userAuth "github.com/gernest/zedlist/middlewares/auth"
	"github.com/gernest/zedlist/middlewares/csrf"
	"github.com/gernest/zedlist/middlewares/flash"
	"github.com/gernest/zedlist/middlewares/i18n"
	"github.com/gernest/zedlist/routes/base"

	"github.com/labstack/echo"
)

var ts = testServer()
var client = &http.Client{Jar: getJar()}
var (
	loginPath    = "/auth/login"
	registerPath = "/auth/register"
	logoutPath   = "/auth/logout"
)

func testServer() *httptest.Server {
	e := echo.New()
	e.Renderer = tmpl.NewRenderer()
	// middlewares
	e.Use(utils.WrapMiddleware(i18n.Langs()))      // languages
	e.Use(utils.WrapMiddleware(flash.Flash()))     // flash messages
	e.Use(utils.WrapMiddleware(userAuth.Normal())) // adding user context data

	// HOME
	e.GET("/", base.Home)

	// AUTH
	xauth := e.Group("/auth")

	xauth.Use(echo.WrapMiddleware(csrf.Nosurf()))  // csrf protection
	xauth.Use(utils.WrapMiddleware(csrf.Tokens())) // csrf tokens

	xauth.GET("/login", Login)
	xauth.POST("/login", LoginPost)
	xauth.GET("/register", Register)
	xauth.POST("/register", RegisterPost)
	xauth.GET("/logout", Logout)

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
//		AUTH ROUTES
//
//
func TestGetLogin(t *testing.T) {
	l := fmt.Sprintf("%s%s", ts.URL, loginPath)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting login page %v", err)
	}

	// The title of the page should be set to login
	title := "<title>login</title>"
	if !bytes.Contains(b, []byte(title)) {
		t.Errorf(" expected login page got %s", b)
	}
}
func TestGetRegister(t *testing.T) {
	l := fmt.Sprintf("%s%s", ts.URL, registerPath)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting login page %v", err)
	}

	// The title of the page should be set to register
	title := "<title>register</title>"
	if !bytes.Contains(b, []byte(title)) {
		t.Errorf(" expected login page got %s", b)
	}
}
func TestPostRegister(t *testing.T) {
	l := fmt.Sprintf("%s%s", ts.URL, registerPath)
	vars := url.Values{
		"first_name":       {"geofrey"},
		"last_name":        {"enrnest"},
		"middle_name":      {"gernest"},
		"email":            {"gernest@mawazo.com"},
		"password":         {"kilimahewa"},
		"confirm_password": {"kilimahewa"},
		"gender":           {"1"},
		"birth_date":       {"2 January, 1980"},
	}

	// lets obtain the csrf_token to submit with the form.
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting login page %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		t.Errorf("loading document %v", err)
	}
	token, ok := doc.Find("#token").Attr("value")
	if !ok {
		t.Errorf("expected crsf to ken to be set")
	}
	vars.Set("csrf_token", token)

	resp, err := client.PostForm(l, vars)
	if err != nil {
		t.Errorf(" posting registration form %v", err)
	}
	defer resp.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)

	// SHould redirect to login page if registration is successful
	// The title of the page should be set to login
	title := "<title>login</title>"
	if !bytes.Contains(buf.Bytes(), []byte(title)) {
		t.Errorf(" expected login page got %s", buf)
	}
}
func TestPostLogin(t *testing.T) {
	l := fmt.Sprintf("%s%s", ts.URL, loginPath)

	// Obtain the csrf token
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting login page %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		t.Errorf("loading document %v", err)
	}
	token, ok := doc.Find("#token").Attr("value")
	if !ok {
		t.Errorf("expected crsf to ken to be set")
	}

	vars := url.Values{
		"email":      {"gernest@mawazo.com"},
		"password":   {"kilimahewa"},
		"csrf_token": {token},
	}

	// case invalid form
	vars.Set("email", "boom")
	resp0, err := client.PostForm(l, vars)
	if err != nil {
		t.Errorf(" posting registration form %v", err)
	}
	defer resp0.Body.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, resp0.Body)
	loginTitle := "<title>login</title>"
	if !bytes.Contains(buf.Bytes(), []byte(loginTitle)) {
		t.Errorf(" expected home page got %s", buf)
	}

	// case user not found
	vars.Set("email", "bigbang@space.com")
	resp1, err := client.PostForm(l, vars)
	if err != nil {
		t.Errorf(" posting registration form %v", err)
	}
	defer resp1.Body.Close()
	buf.Reset()
	io.Copy(buf, resp1.Body)
	if !bytes.Contains(buf.Bytes(), []byte(loginTitle)) {
		t.Errorf(" expected home page got %s", buf)
	}

	// case a passing login
	vars.Set("email", "gernest@mawazo.com") // restore the field
	resp, err := client.PostForm(l, vars)
	if err != nil {
		t.Errorf(" posting registration form %v", err)
	}
	defer resp.Body.Close()

	buf.Reset()
	io.Copy(buf, resp.Body)

	// SHould redirect to login page if registration is successful
	// The title of the page should be set to login
	title := "<title>zedlist</title>"
	if !bytes.Contains(buf.Bytes(), []byte(title)) {
		t.Errorf(" expected home page got %s", buf)
	}

	// should contain the logout button
	outButton := "logout"
	if !bytes.Contains(buf.Bytes(), []byte(outButton)) {
		t.Errorf(" expected home page with logout button  got %s", buf)
	}
}

func TestLogout(t *testing.T) {
	l := fmt.Sprintf("%s%s", ts.URL, logoutPath)
	b, err := com.HttpGetBytes(client, l, nil)
	if err != nil {
		t.Errorf("getting login page %v", err)
	}

	// should redirect to home page
	title := "<title>zedlist</title>"
	if !bytes.Contains(b, []byte(title)) {
		t.Errorf(" expected home page got %s", b)
	}
	// should not contain the logout button
	outButton := "logout"
	if bytes.Contains(b, []byte(outButton)) {
		t.Errorf(" expected home page without logout button  got %s", b)
	}
}

// this should be always at the bottom of this file.
func TestCLose(t *testing.T) {
	closeTest()
}
