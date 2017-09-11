// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package i18n

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/gernest/zedlist/modules/utils"

	"github.com/labstack/echo"
)

func TestLangs(t *testing.T) {
	e := echo.New()
	e.Use(utils.WrapMiddleware(Langs()))
	e.GET("/", func(ctx echo.Context) error {
		lang := utils.GetLang(ctx)
		return ctx.String(http.StatusOK, lang)
	})
	ts := httptest.NewServer(e)
	defer ts.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{Jar: jar}

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)

	// The default language should be en.
	if buf.String() != "en" {
		t.Errorf("expected en got %s", buf.String())
	}

}
