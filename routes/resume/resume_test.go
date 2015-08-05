// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package resume

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"

	"github.com/labstack/echo"
)

var ts = testServer()
var user = testUser()

func userMiddleware(ctx *echo.Context) error {
	ctx.Set("User", user)
	return nil
}

func testServer() *echo.Echo {
	e := echo.New()
	e.SetRenderer(tmpl.NewRenderer())
	e.Use(userMiddleware)
	e.Get("/resume/", Home)
	e.Post("/resume/create", Create)
	e.Get("/resume/update/:id", Update)
	e.Get("/resume/view/:id", View)
	e.Post("/resume/delete", Delete)
	return e
}
func testUser() *models.Person {
	p := &models.Person{
		AboutMe: "nut cracker",
	}
	query.Create(p)
	return p
}

func TestHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/resume/", nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.Code)
	}
}

func TestCreate(t *testing.T) {
	vars := url.Values{
		"resume_name": {"red dragon"},
	}
	req, err := http.NewRequest("POST", "/resume/create", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	req.Header.Set(echo.ContentType, echo.ApplicationForm)
	resp := httptest.NewRecorder()
	ctx := echo.NewContext(req, echo.NewResponse(resp), ts)
	ctx.Set("User", user)

	err = Create(ctx)
	if err != nil {
		t.Error(err)
	}
	if resp.Code != http.StatusFound {
		t.Errorf("expected  %d got %d", http.StatusFound, resp.Code)
	}
}

func TestView(t *testing.T) {
	path := fmt.Sprintf("/resume/view/%d", user.Resumes[0].ID)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.Code)
	}
}
