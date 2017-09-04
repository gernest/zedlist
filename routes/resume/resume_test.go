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

	"github.com/gernest/zedlist/middlewares/i18n"
	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/labstack/echo"
)

var ts = testServer()
var user = testUser()

func userMiddleware(ctx echo.Context) error {
	ctx.Set("User", user)
	return nil
}

func testServer() *echo.Echo {
	e := echo.New()
	e.Renderer = tmpl.NewRenderer()
	e.Use(utils.WrapMiddleware(i18n.Langs())) // languages
	e.Use(utils.WrapMiddleware(userMiddleware))
	e.GET("/resume/", Home)
	e.POST("/resume/create", Create)
	e.GET("/resume/update/:id", Update)
	e.GET("/resume/view/:id", View)
	e.POST("/resume/delete/:id", Delete)
	return e
}
func testUser() *models.Person {
	p := &models.Person{
		AboutMe: "nut cracker",
	}
	query.Create(db.Conn, p)
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
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
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

func TestUpdate(t *testing.T) {
	path := fmt.Sprintf("/resume/update/%d", user.Resumes[0].ID)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.Code)
	}

	// Case the resume is not found
	path = "/resume/update/200"
	req, err = http.NewRequest("GET", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp = httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected %d got %d", http.StatusNotFound, resp.Code)
	}

	// Case bad request
	path = "/resume/update/mia"
	req, err = http.NewRequest("GET", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp = httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, resp.Code)
	}
}

func TestDelete(t *testing.T) {
	path := fmt.Sprintf("/resume/delete/%d", user.Resumes[0].ID)
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusFound {
		t.Errorf("expected %d got %d", http.StatusFound, resp.Code)
	}

	// Case the resume is not found
	path = "/resume/delete/200"
	req, err = http.NewRequest("POST", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp = httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected %d got %d", http.StatusNotFound, resp.Code)
	}

	// Case bad request
	path = "/resume/delete/mia"
	req, err = http.NewRequest("POST", path, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	resp = httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, resp.Code)
	}
}
