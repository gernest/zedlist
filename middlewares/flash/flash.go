// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package flash is a collection of  flash messages middleware for zedlist.
package flash

import (
	"github.com/labstack/echo"
	"github.com/zedio/zedlist/modules/flash"
)

// Flash adds flash messages to the request context.
func Flash() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		flash.AddFlashToCtx(ctx)
		return nil
	}
}
