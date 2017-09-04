// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package dash is a collection of user dashboard handlers for zedlist.
package dash

import (
	"net/http"

	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/flash"

	"github.com/gernest/zedlist/modules/forms"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"

	"github.com/gorilla/schema"
	"github.com/labstack/echo"
)

var msgInvalidorm = "some fish happende"

var formDecoder = schema.NewDecoder()

// Home renders dashboard home page.
//
//		Method           GET
//
//		Route            /dash/
//
//		Restrictions     Yes
//
// 		Template         dash/home.html
func Home(ctx echo.Context) error {
	utils.SetData(ctx, "PageTitle", "dashboard")
	f := forms.New(utils.GetLang(ctx))
	utils.SetData(ctx, "JobForm", f.JobForm()())
	return ctx.Render(http.StatusOK, tmpl.DashHomeTpl, utils.GetData(ctx))
}

// JobsNewGet renders the new job form.
//
//		Method           GET
//
//		Route            /dash/jobs/new
//
//		Restrictions     Yes
//
// 		Template         dash/jobs_new.html
func JobsNewGet(ctx echo.Context) error {
	f := forms.New(utils.GetLang(ctx))
	utils.SetData(ctx, "PageTitle", "new job")
	utils.SetData(ctx, "JobForm", f.JobForm()())
	return ctx.Render(http.StatusOK, tmpl.DashJobTpl, utils.GetData(ctx))
}

// JobsNewPost process the new job form.
//
//		Method           POST
//
//		Route            /dash/jobs/new
//
//		Restrictions     Yes
//
// 		Template         None
func JobsNewPost(ctx echo.Context) error {
	var flashMessages = flash.New()
	f := forms.New(utils.GetLang(ctx))
	jf := f.JobForm()(ctx.Request())
	if !jf.IsValid() {
		// TODO: improve flash message ?
		flashMessages.Err(msgInvalidorm)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/dash/jobs/new")
		return nil
	}

	if isLoged := ctx.Get("IsLoged"); isLoged != nil {
		person := ctx.Get("User").(*models.Person)
		if jerr := query.PersonCreateJob(db.Conn, person, jf.GetModel().(forms.JobForm)); jerr != nil {
			// TODO: improve flash message ?
			flashMessages.Err("some really bad fish happened")
			flashMessages.Save(ctx)
			ctx.Redirect(http.StatusFound, "/dash/jobs/new")
			return nil
		}
		// add flash message
		flashMessages.Success("new job was created successful")
		flashMessages.Save(ctx)

		ctx.Redirect(http.StatusFound, "/dash/")
		return nil
	}
	he := echo.NewHTTPError(http.StatusUnauthorized)
	ctx.Error(he)
	return nil
}

// Profile renders user profile.
//
//		Method           GET
//
//		Route            /dash/profile
//
//		Restrictions     Yes
//
// 		Template         dash/profile.html
func Profile(ctx echo.Context) error {
	utils.SetData(ctx, "PageTitle", "profile")
	return ctx.Render(http.StatusOK, tmpl.DashProfileTpl, utils.GetData(ctx))
}

// ProfileName updates Person's names.
//
//		Method       POST
//
//		Route         /dash/profile/name
//
//		Restricted    Yes
//
//		Template      None (everything is redirected to '/dash/profile' )
//
// When there are validation errors flash messages are set.
func ProfileName(ctx echo.Context) error {
	r := ctx.Request()
	v := forms.NewValid(utils.GetLang(ctx))
	r.ParseForm()
	pName := &models.PersonName{}
	if err := formDecoder.Decode(pName, r.PostForm); err != nil {
		// TODO: do something?
	}
	errs := v.ValidatePersonName(pName)
	if errs != nil {
		// TODO: do somethins?
	}

	person := ctx.Get("User").(*models.Person)
	person.UpdateNames(pName)
	err := query.Update(db.Conn, person)
	if err != nil {
		// TODO: do somethins?
	}

	ctx.Redirect(http.StatusFound, "/dash/profile")
	return nil
}
