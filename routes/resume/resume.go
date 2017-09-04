// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package resume has resumes handlers for zedlist.
package resume

import (
	"fmt"
	"net/http"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/flash"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/labstack/echo"
)

// Home shows the resumes home page.
//
//		Method           GET
//
//		Route            /resume/
//
//		Restrictions     Yes
//
// 		Template         resume/home.html
func Home(ctx echo.Context) error {
	user := ctx.Get("User").(*models.Person)
	if res, err := query.GetAllPersonResumes(db.Conn, user); err == nil {
		utils.SetData(ctx, "resumes", res)
	}
	return ctx.Render(http.StatusOK, tmpl.ResumeHomeTpl, utils.GetData(ctx))
}

// View displays the resume.
//
//		Method           GET
//
//		Route            /resume/view
//
//		Restrictions     Yes
//
// 		Template         resume/view.html
func View(ctx echo.Context) error {
	iid, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	resume, err := query.GetResumeByID(db.Conn, iid)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, tmpl.NotFoundMessage)
	}
	utils.SetData(ctx, "resume", resume)
	return ctx.Render(http.StatusOK, tmpl.ResumeViewTpl, utils.GetData(ctx))
}

// Create creates a new resume.
//
//		Method           POST
//
//		Route            /resume/new
//
//		Restrictions     Yes
//
// 		Template         None
func Create(ctx echo.Context) error {
	var flashMessages = flash.New()
	r := ctx.Request()
	r.ParseForm()
	name := r.Form.Get("resume_name")

	user := ctx.Get("User").(*models.Person)

	resume := models.SampleResume()
	resume.Name = name
	err := query.CreateResume(db.Conn, user, resume)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.ServerErrorMessage)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}

	flashMessages.Success("successful created a new resume")
	flashMessages.Save(ctx)

	// Redirect to the update page for further updating of the resume.
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/resume/update/%d", resume.ID))
	return nil
}

// Update renders the resume update page.
//
//		Method           GET
//
//		Route            /resume/update/:id
//
//		Restrictions     Yes
//
// 		Template         None
func Update(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	user := ctx.Get("User").(*models.Person)

	resume, err := query.GetResumeByID(db.Conn, id)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}

	// Users are allowed to update resumes that they own.
	if resume.PersonID != user.ID {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	utils.SetData(ctx, "resume", resume)
	return ctx.Render(http.StatusOK, tmpl.ResumeUpddateTpl, utils.GetData(ctx))
}

// Delete deletes the resume.
//
//		Method           POST
//
//		Route            /resume/delete/:id
//
//		Restrictions     Yes
//
// 		Template         None
func Delete(ctx echo.Context) error {
	var flashMessages = flash.New()
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	user := ctx.Get("User").(*models.Person)

	resume, err := query.GetResumeByID(db.Conn, id)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}

	// Users are allowed to delete resumes that they don't own.
	if resume.PersonID != user.ID {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}

	err = query.Delete(db.Conn, resume)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.ServerErrorMessage)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}

	flashMessages.Success("resume successful deleted")
	flashMessages.Save(ctx)
	ctx.Redirect(http.StatusFound, "/resume/")
	return nil
}
