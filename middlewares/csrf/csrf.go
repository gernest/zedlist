// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package csrf is a collection of csrf protection middlewares for zedlist.
package csrf

import (
	"net/http"

	"github.com/gernest/zedlist/modules/utils"

	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
)

// Nosurf is a wrapper for justinas' csrf protection middleware
func Nosurf() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nosurf.New(next)
	}
}

// Tokens adds csrf token context. the context key is CsrfToken
func Tokens() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		utils.SetData(ctx, "CsrfToken", nosurf.Token(ctx.Request()))
		return nil
	}
}
