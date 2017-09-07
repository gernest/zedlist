// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package i18n contains translation middleware for zedlist.
package i18n

import (
	"github.com/gernest/zedlist/modules/utils"

	"github.com/gernest/zedlist/modules/session"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/labstack/echo"
)

var store = session.New()

type Lang struct {
	Long  string
	Short string
	Flag  string
}

func SupportedLangs() []Lang {
	return []Lang{
		{
			Long:  "Swahili",
			Short: "sw",
			Flag:  "tz",
		},
		{
			Long:  "English",
			Short: "en",
			Flag:  "us",
		},
	}
}

// Langs sets active language in the request context.
func Langs() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		utils.SetData(ctx, settings.SupportedLangs, SupportedLangs())
		sess, _ := store.Get(ctx.Request(), settings.App.Session.Lang)
		target := sess.Values[settings.LangSessionKey]
		if target != nil {
			utils.SetData(ctx, settings.LangDataKey, target)
			return nil
		}
		sess.Values[settings.LangDataKey] = settings.App.DefaultLang
		store.Save(ctx.Request(), ctx.Response(), sess)
		utils.SetData(ctx, settings.LangDataKey, settings.App.DefaultLang)
		return nil
	}
}
