// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package flash

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zedio/zedlist/modules/flash"
	"github.com/zedio/zedlist/modules/utils"

	"github.com/kr/pretty"
	"github.com/labstack/echo"
)

var message = "hello flash"

func helloFlash(ctx echo.Context) error {
	flashMessages := flash.New()
	flashMessages.Success(message)
	flashMessages.Save(ctx)
	ctx.Redirect(http.StatusFound, "/home")
	return nil
}

func home(ctx echo.Context) error {
	d := utils.GetData(ctx).(utils.Data)
	flashes := d.Get("Flash").(flash.Flashes)[0]
	return ctx.String(http.StatusOK, fmt.Sprintf("%#v", pretty.Formatter(flashes)))
}

func TestFlash(t *testing.T) {
	e := echo.New()
	e.Use(utils.WrapMiddleware(Flash()))
	e.GET("/home", home)
	e.GET("/flash", helloFlash)
	ts := httptest.NewServer(e)
	defer ts.Close()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{Jar: jar}

	resp, err := client.Get(fmt.Sprintf("%s/flash", ts.URL))
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, resp.StatusCode)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	if !strings.Contains(buf.String(), message) {
		t.Errorf("expected %s to contain %s", buf.String(), message)
	}

	// resp1, err := client.Get(fmt.Sprintf("%s/home", ts.URL))
	// if err != nil {
	// 	t.Error(err)
	// }
	// defer resp1.Body.Close()
	// buf.Reset()
	// io.Copy(buf, resp1.Body)
	// if !strings.Contains(buf.String(), message) {
	// 	t.Errorf("expected %s to contain %s", buf.String(), message)
	// }

}
