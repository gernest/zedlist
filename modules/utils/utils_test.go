// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gernest/zedlist/modules/settings"

	"github.com/labstack/echo"
)

func TestData(t *testing.T) {
	sampleData := []struct {
		key string
		val interface{}
	}{
		{"one", 1},
		{"two", "two"},
		{"three", "three"},
	}
	e := echo.New()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()

	// Create the context on which we will be messing with
	ctx := echo.NewContext(req, echo.NewResponse(res), e)

	//
	//	SetData
	//
	for _, v := range sampleData {
		SetData(ctx, v.key, v.val)
	}

	//
	//	GetData
	//
	data, ok := GetData(ctx).(Data)
	if !ok {
		t.Error("kaboom")
	}

	// Check the Data methods
	d := data.Get(sampleData[0].key)
	if d == nil {
		t.Errorf("expected %v got nil instead", sampleData[0].val)
	}

	data.Set("hello", "world")
	if h := data.Get("hello"); h != nil {
		if h.(string) != "world" {
			t.Errorf("expected world got %v", h)
		}
	}

	//
	//	GetLang
	//
	SetData(ctx, settings.LangDataKey, "en")

	lang := GetLang(ctx)
	if lang != "en" {
		t.Errorf("expected en got %s", lang)
	}
}

func TestIsAjax(t *testing.T) {
	e := echo.New()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	ctx := echo.NewContext(req, echo.NewResponse(res), e)

	if ok := IsAjax(ctx); ok {
		t.Errorf("expected false got %v", ok)
	}

	req1, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	req1.Header.Set("X-Requested-With", "XMLHttpRequest")
	res1 := httptest.NewRecorder()
	ctx1 := echo.NewContext(req1, echo.NewResponse(res1), e)

	if ok := IsAjax(ctx1); !ok {
		t.Errorf("expected true got %v", ok)
	}
}

func TestGetInt(t *testing.T) {
	one, err := GetInt("1")
	if err != nil {
		t.Error(err)
	}
	if one != 1 {
		t.Errorf("expected 1 got %v", one)
	}
}
