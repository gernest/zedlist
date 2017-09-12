// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package resume has resumes handlers for zedlist.
package resume

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/zedio/zedlist/models"
	"github.com/zedio/zedlist/modules/db"
	"github.com/zedio/zedlist/modules/flash"
	"github.com/zedio/zedlist/modules/log"
	"github.com/zedio/zedlist/modules/query"
	"github.com/zedio/zedlist/modules/tmpl"
	"github.com/zedio/zedlist/modules/utils"

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

func New(ctx echo.Context) error {
	utils.SetData(ctx, "PageTitle", "new resume")
	utils.SetData(ctx, "Scripts", []template.HTML{
		template.HTML(`/static/js/moon.min.js`),
		template.HTML(`/static/js/monx.min.js`),
		template.HTML(`/static/js/resume.js`),
	})
	return ctx.Render(http.StatusOK, tmpl.ResumeNewTpl, utils.GetData(ctx))
}

type resumeReq struct {
	ID        int64  `json:"id,omitempty"`
	ProfileID int64  `json:"profileID"`
	Title     string `json:"title"`
}

func Update(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	req := &resumeReq{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, models.NewJSONErr(http.StatusText(
			http.StatusUnprocessableEntity,
		)))
	}
	rs, err := query.GetResumeByID(db.Conn, req.ID)
	if err != nil {
		log.Error(ctx, err)
		if query.NotFound(err) {
			return ctx.JSON(http.StatusNotFound, models.NewJSONErr(http.StatusText(
				http.StatusNotFound,
			)))
		}
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSON(http.StatusOK, rs)
}

func NewPost(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	req := &resumeReq{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, models.NewJSONErr(http.StatusText(
			http.StatusUnprocessableEntity,
		)))
	}
	rs := &models.Resume{
		PersonID: req.ProfileID,
		Title:    req.Title,
	}
	if err = query.Create(db.Conn, rs); err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSON(http.StatusOK, rs)
}

func BasicPut(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	b := &models.Basic{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, b)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, models.NewJSONErr(http.StatusText(
			http.StatusUnprocessableEntity,
		)))
	}
	if err = query.Update(db.Conn, b); err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusOK, b, "\t")
}

func BasicPost(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		ctx.Echo().DefaultHTTPErrorHandler(echo.ErrUnsupportedMediaType, ctx)
		return nil
	}
	b := &models.Basic{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, b)
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	if err = query.Create(db.Conn, b); err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusCreated, b, "\t")
}

func BasicGet(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	b, err := query.GetBasicResumeByID(db.Conn, id)
	if err != nil {
		log.Error(ctx, err)
		if query.NotFound(err) {
			return ctx.JSON(http.StatusNotFound, models.NewJSONErr(http.StatusText(
				http.StatusNotFound,
			)))
		}
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusOK, b, "\t")
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
	resume.Title = name
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
