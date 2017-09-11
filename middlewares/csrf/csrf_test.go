// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package csrf

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/zedio/zedlist/modules/utils"

	"github.com/labstack/echo"
)

func hello(ctx echo.Context) error {
	token := ""
	d := utils.GetData(ctx).(utils.Data)
	tok := d.Get("CsrfToken")
	if tok != nil {
		token = tok.(string)
	}
	return ctx.String(http.StatusOK, token)
}

func TestCsrf(t *testing.T) {
	e := echo.New()
	e.Use(echo.WrapMiddleware(Nosurf()))
	e.Use(utils.WrapMiddleware(Tokens()))
	e.POST("/", hello)
	e.GET("/", hello)

	ts := httptest.NewServer(e)
	defer ts.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Log
	}
	client := &http.Client{Jar: jar}
	formVars := url.Values{
		"nothing": {"to show"},
	}

	res, err := client.PostForm(ts.URL, formVars)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, res.StatusCode)
	}

	// get the token
	res1, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer res1.Body.Close()

	b := &bytes.Buffer{}
	io.Copy(b, res1.Body)

	if res1.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, res1.StatusCode)
	}
	formVars.Set("csrf_token", b.String())

	// now we a full qualified
	res2, err := client.PostForm(ts.URL, formVars)
	if err != nil {
		t.Error(err)
	}
	defer res2.Body.Close()

	if res2.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, res2.StatusCode)
	}
}
