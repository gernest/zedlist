// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package utils contains helpers for zedlist.
package utils

import (
	"github.com/Unknwon/com"
	"github.com/labstack/echo"
	"github.com/zedio/zedlist/modules/session"
	"github.com/zedio/zedlist/modules/settings"
)

var store = session.New()

// Data holds templates context data.
type Data map[string]interface{}

// Get retrueves a value by key
func (d Data) Get(key string) interface{} {
	return d[key]
}

// Set sets key value to val.
func (d Data) Set(key string, val interface{}) {
	d[key] = val
}

// GetContext returns a Data object with default values.
func GetContext() Data {
	d := make(Data)
	d["AppName"] = settings.App.Name
	d["AppUrl"] = settings.App.AppURL
	return d
}

// SetData stores the given key value in the *echo.Context(under the template context oject)
func SetData(ctx echo.Context, key string, val interface{}) {
	v := ctx.Get(settings.DataKey)
	switch v.(type) {
	case Data:
		d := v.(Data)
		d[key] = val
		ctx.Set(settings.DataKey, d)
	default:
		d := GetContext()
		d[key] = val
		ctx.Set(settings.DataKey, d)
	}

}

// GetData returns template context data stored in *echo.Context
func GetData(ctx echo.Context) interface{} {
	data := GetContext()
	if v := ctx.Get(settings.DataKey); v != nil {
		data = v.(Data)
	}
	if v := ctx.Get(settings.FlashKey); v != nil {
		data.Set(settings.FlashKey, v)
	}
	if v := ctx.Get(settings.FlashCtxKey); v != nil {
		data.Set(settings.FlashCtxKey, v)
	}
	return data
}

// GetLang retrieves language from the context.
func GetLang(ctx echo.Context) string {
	return GetData(ctx).(Data).Get(settings.LangDataKey).(string)
}

// IsAjax returns true if the request is ajax and false otherwose..
func IsAjax(ctx echo.Context) bool {
	return ctx.Request().Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// DeleteSession delete session by name.
func DeleteSession(ctx echo.Context, sessionName string) error {
	ss, err := store.Get(ctx.Request(), sessionName)
	if err != nil {
		return err
	}
	return store.Delete(ctx.Request(), ctx.Response(), ss)
}

//GetInt converts the string to int
func GetInt64(str string) (int64, error) {
	return com.StrTo(str).Int64()
}

//GetInt converts the string to int
func GetInt(str string) (int, error) {
	return com.StrTo(str).Int()
}

// WrapMiddleware wraps a echo.HandlerFunc to echo.MiddlewareFunc
func WrapMiddleware(h echo.HandlerFunc) echo.MiddlewareFunc {
	return func(a echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := h(ctx); err != nil {
				return err
			}
			return a(ctx)
		}
	}
}
