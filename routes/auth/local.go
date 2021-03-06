// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package auth is a collection of auhentication handlers for zedlist.
package auth

import (
	"net/http"

	validate "github.com/asaskevich/govalidator"
	"github.com/zedio/zedlist/models"
	"github.com/zedio/zedlist/modules/db"
	"github.com/zedio/zedlist/modules/flash"
	"github.com/zedio/zedlist/modules/forms"
	"github.com/zedio/zedlist/modules/log"
	"github.com/zedio/zedlist/modules/query"
	"github.com/zedio/zedlist/modules/utils"

	"github.com/labstack/echo"
	"github.com/zedio/zedlist/modules/session"
	"github.com/zedio/zedlist/modules/settings"
	"github.com/zedio/zedlist/modules/tmpl"
)

var sessStore = session.New()

// Login renders login form.
//
//		Method           GET
//
//		Route            /auth/login
//
//		Restrictions     None
//
// 		Template         auth/login.html
//
func Login(ctx echo.Context) error {

	f := forms.New(utils.GetLang(ctx))
	utils.SetData(ctx, "form", f)

	// set page tittle to login
	utils.SetData(ctx, settings.PageTitle, "login")

	return ctx.Render(http.StatusOK, tmpl.LoginTpl, utils.GetData(ctx))
}

// LoginPost handlers login form, and logs in the user. If the form is valid,
// the user is redirected to "/auth/login" with the form validation errors. When
// the user is validated redirection is made to "/".
//
//		Method           POST
//
//		Route            /auth/login
//
//		Restrictions     None
//
// 		Template         None (All actions redirect to other routes )
//
// Flash messages may be set before redirection.
func LoginPost(ctx echo.Context) error {
	var flashMessages = flash.New()
	f := forms.New(utils.GetLang(ctx))
	lf, err := f.DecodeLogin(ctx.Request())
	if err != nil {
		ctx.Redirect(http.StatusFound, "/auth/login")
		return nil
	}
	if !lf.Valid() {
		for k, v := range lf.Ctx() {
			flashMessages.AddCtx(k, v)
		}
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/auth/login")
		return nil
	}
	var user *models.User
	if validate.IsEmail(lf.Name) {
		user, err = query.AuthenticateUserByEmail(db.Conn, *lf)
		if err != nil {
			log.Error(ctx, err)

			// We want the user to try again, but rather than rendering the form right
			// away, we redirect him/her to /auth/login route(where the login process with
			// start aflsesh albeit with a flash message)
			flashMessages.Err(settings.FlashLoginErr)
			flashMessages.Save(ctx)
			ctx.Redirect(http.StatusFound, "/auth/login")
			return nil
		}
	} else {
		user, err = query.AuthenticateUserByName(db.Conn, *lf)
		if err != nil {
			log.Error(ctx, err)

			// We want the user to try again, but rather than rendering the form right
			// away, we redirect him/her to /auth/login route(where the login process with
			// start aflsesh albeit with a flash message)
			flashMessages.Err(settings.FlashLoginErr)
			flashMessages.Save(ctx)
			ctx.Redirect(http.StatusFound, "/auth/login")
			return nil
		}
	}

	// create a session for the user after the validation has passed. The info stored
	// in the session is the user ID, where as the key is userID.
	ss, err := sessStore.Get(ctx.Request(), settings.App.Session.Name)
	if err != nil {
		log.Error(ctx, err)
	}
	ss.Values["userID"] = user.ID
	err = ss.Save(ctx.Request(), ctx.Response())
	if err != nil {
		log.Error(ctx, err)
	}
	person, err := query.GetPersonByUserID(db.Conn, user.ID)
	if err != nil {
		log.Error(ctx, err)
		flashMessages.Err(settings.FlashLoginErr)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/auth/login")
		return nil
	}

	// add context data. IsLoged is just a conveniece in template rendering. the User
	// contains a models.Person object, where the PersonName is already loaded.
	utils.SetData(ctx, "IsLoged", true)
	utils.SetData(ctx, "User", person)
	flashMessages.Success(settings.FlashLoginSuccess)
	flashMessages.Save(ctx)
	ctx.Redirect(http.StatusFound, "/")
	return nil
}

// Register renders registration form.
//
//		Method           GET
//
//		Route            /auth/register
//
//		Restrictions     None
//
// 		Template         auth/register.html
func Register(ctx echo.Context) error {
	// set page tittle to register
	utils.SetData(ctx, "PageTitle", "register")
	return ctx.Render(http.StatusOK, tmpl.RegisterTpl, utils.GetData(ctx))
}

// RegisterPost handles registration form, and create a session for the new user if the registration
// process is complete.
//
//		Method           POST
//
//		Route            /auth/register
//
//		Restrictions     None
//
// 		Template         None (All actions redirect to other routes )
//
// Flash messages may be set before redirection.
func RegisterPost(ctx echo.Context) error {
	var flashMessages = flash.New()
	f := forms.New(utils.GetLang(ctx))

	lf, err := f.DecodeRegister(ctx.Request())
	if err != nil {
		// Case the form is not valid, ships it back with the errors exclusively
		return ctx.Render(http.StatusOK, tmpl.RegisterTpl, utils.GetData(ctx))
	}
	if !lf.Valid() {
		for k, v := range lf.Ctx() {
			flashMessages.AddCtx(k, v)
		}
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, "/auth/register")
	}
	// we are not interested in the returned user, rather we make sure the user has
	// been created.
	_, err = query.CreateNewUser(db.Conn, *lf)
	if err != nil {
		flashMessages.Err(settings.FlashAccountCreateFailed)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/auth/register")
		return nil
	}

	// TODO: improve the message to include directions to use the current email and
	// password to login?
	flashMessages.Success(settings.FlashAccountCreate)
	flashMessages.Save(ctx)

	// Don't create session in this route, its best to leave only one place which
	// messes with the main user session. So we redirect to the login page, and encourage
	// the user to login.
	ctx.Redirect(http.StatusFound, "/auth/login")
	return nil
}

// Logout deletes all sessions,then redirects to "/".
//
//		Method           POST
//
//		Route            /auth/logout
//
//		Restrictions     None
//
// 		Template         None (All actions redirect to other routes )
//
// Flash messages may be set before redirection.
func Logout(ctx echo.Context) error {
	if _, ok := ctx.Get("User").(*models.Person); ok {
		utils.DeleteSession(ctx, settings.App.Session.Lang)
		utils.DeleteSession(ctx, settings.App.Session.Flash)
		utils.DeleteSession(ctx, settings.App.Session.Name)
	}
	ctx.Redirect(http.StatusFound, "/")
	return nil
}

func Delete(ctx echo.Context) error {
	f := forms.New(utils.GetLang(ctx))
	utils.SetData(ctx, "form", f)

	// set page tittle to login
	utils.SetData(ctx, settings.PageTitle, "confirm deleting account")
	return ctx.Render(http.StatusOK, tmpl.DeleteTpl, utils.GetData(ctx))
}

func DeletePost(ctx echo.Context) error {
	flashMessages := flash.New()
	f := forms.New(utils.GetLang(ctx))
	u, ok := ctx.Get("User").(*models.Person)
	if !ok {
		flashMessages.Err(settings.FlashNotAuthorized)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/auth/delete")
		return nil
	}
	usr := f.DecodeDelete(ctx.Request())
	if usr != u.PersonName.Name {
		flashMessages.Err(settings.FlashUnknownAccount)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/auth/delete")
		return nil
	}
	id := ctx.Get("UserID").(int64)
	utils.DeleteSession(ctx, settings.App.Session.Lang)
	utils.DeleteSession(ctx, settings.App.Session.Flash)
	utils.DeleteSession(ctx, settings.App.Session.Name)

	if err := query.DeleteUser(db.Conn, id); err != nil {
		log.Error(ctx, err)
	}
	return ctx.Redirect(http.StatusFound, "/")
}
