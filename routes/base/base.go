// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package base is a collection of base handlers for zedlist.
package base

import (
	"net/http"
	"net/url"

	"github.com/zedio/zedlist/modules/db"
	"github.com/zedio/zedlist/modules/query"
	"github.com/zedio/zedlist/modules/session"
	"github.com/zedio/zedlist/modules/settings"
	"github.com/zedio/zedlist/modules/tmpl"
	"github.com/zedio/zedlist/modules/utils"

	"github.com/labstack/echo"
)

//
//
//		BASE
//
//

// Home renders zedlist home page.
//
//		Method           GET
//
//		Route            /
//
//		Restrictions     None
//
// 		Template         base/home.html
func Home(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, tmpl.BaseHomeTpl, utils.GetData(ctx))
}

// SetLanguage switches language, between swahili and english. This means when this
// handler is accessed, if the current language is en it will be set to sw and vice versa.
//
//		Method           GET
//
//		Route            /language/:lang
//
//		Restrictions     None
//
// 		Template         None ( Redirection is made to home route ("/"))
func SetLanguage(ctx echo.Context) error {
	lang := ctx.Param("lang")
	p := ctx.QueryParam("path")
	q, _ := url.QueryUnescape(p)
	store := session.New()
	sess, _ := store.Get(ctx.Request(), settings.LangSessionName)
	sess.Values[settings.LangSessionKey] = lang
	store.Save(ctx.Request(), ctx.Response(), sess)
	return ctx.Redirect(http.StatusFound, q)
}

//
//
//		REGIONS
//
//

// RegionsHome renders regions home page.
//
//
//		Method           GET
//
//		Route            /jobs/regions
//
//		Restrictions     None
//
// 		Template         base/regions.html
//
func RegionsHome(ctx echo.Context) error {
	regs, err := query.GetAllRegions(db.Conn)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, settings.RegionsListKey, regs)
	return ctx.Render(http.StatusOK, tmpl.BaseRegionsTpl, utils.GetData(ctx))
}

// RegionsJobView renders jobs from a gien region. The region name should be in short form.
//
//
//		Method           GET
//
//		Route            /jobs/regions/:name
//
//		Restrictions     None
//
// 		Template         base/regions_job.html
//
func RegionsJobView(ctx echo.Context) error {
	name := ctx.Param("name")
	jobs, count, err := query.GetJobByRegionShort(db.Conn, name)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, settings.JobsFound, count)
	utils.SetData(ctx, settings.JobsListKey, jobs)
	return ctx.Render(http.StatusOK, tmpl.BaseRegionsJobViewTpl, utils.GetData(ctx))
}

// RegionsJobPaginate a route frr /jobs/regions/:name/:from/:to. It handles pagination where
// form to is offset and limit respectively.
//
// For example route "/jobs/regions/mza/2/4" will render from 2nd to 4th latest jobs from mwanza.
//
//		Method           GET
//
//		Route            /jobs/regions/:name/:from/:to
//
//		Restrictions     None
//
// 		Template         base/regions.html
func RegionsJobPaginate(ctx echo.Context) error {
	name := ctx.Param("name")
	offset, err := utils.GetInt(ctx.Param("from"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}

	limit, err := utils.GetInt(ctx.Param("to"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	jobs, err := query.GetJobByRegionPaginate(db.Conn, name, offset, limit)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, settings.JobsListKey, jobs)
	return ctx.Render(http.StatusOK, tmpl.BaseRegionsPaginateTpl, utils.GetData(ctx))
}
