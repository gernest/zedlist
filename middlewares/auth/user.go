// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package auth is a collection of authentication middlewares for zedlist.
package auth

import (
	"net/http"

	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/session"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/gernest/zedlist/modules/settings"
	"github.com/labstack/echo"
)

var store = session.New()

// Normal just adds user detail to templates data context. it is wise to add this before
// Must(). This unlike Must will not return an error, regardless of if the user is loged
// in or not
func Normal() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if v := ctx.Get("IsLoged"); v != nil && v.(bool) == true {
			return nil
		}
		ss, err := store.Get(ctx.Request(), settings.App.Session.Name)
		if err != nil {
			log.Error(ctx, err)
		}
		if v, ok := ss.Values["userID"]; ok {
			id := v.(int64)
			person, err := query.GetPersonByUserID(db.Conn, id)
			if err != nil {
				log.Error(ctx, err)
				return err
			}
			ctx.Set("IsLoged", true)
			ctx.Set("User", person)
			ctx.Set("UserID", id)
			utils.SetData(ctx, "IsLoged", true)
			utils.SetData(ctx, "Person", person)
			utils.SetData(ctx, "UserID", id)
			return nil
		}
		return nil
	}
}

// Must ensures that any route is authorized to access the next handler
// otherwise an error is returned.
//
// TODO custom not authorized handler?
func Must() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if v := ctx.Get("IsLoged"); v != nil && v.(bool) == true {
			return nil
		}
		ss, err := store.Get(ctx.Request(), settings.App.Session.Name)
		if err != nil {
			log.Error(ctx, err)
		}
		if v, ok := ss.Values["userID"]; ok {
			id := v.(int64)
			person, err := query.GetPersonByUserID(db.Conn, id)
			if err != nil {
				log.Error(ctx, err)
				return err
			}
			ctx.Set("IsLoged", true)
			ctx.Set("User", person)
			ctx.Set("UserID", id)
			utils.SetData(ctx, "IsLoged", true)
			utils.SetData(ctx, "User", person)
			utils.SetData(ctx, "UserID", id)
			return nil
		}
		return ctx.Redirect(http.StatusFound, "/auth/login")
	}
}
