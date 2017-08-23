// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package base is a collection of base handlers for zedlist.
package base

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gernest/zedlist/bindata/static"
	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/session"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

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

// JobsHome renders jobs home page
//
//
//		Method           GET
//
//		Route            /jobs/
//
//		Restrictions     None
//
// 		Template         base/jobs.html
func JobsHome(ctx echo.Context) error {
	jobs, err := query.GetLatestJobs()
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))

	}
	utils.SetData(ctx, settings.JobsListKey, jobs)
	utils.SetData(ctx, settings.PageTitleKey, "jobs")
	return ctx.Render(http.StatusOK, tmpl.BaseJobsHomeTpl, utils.GetData(ctx))
}

// JobView displays a single job by the given job id.
//
//
//		Method           GET
//
//		Route            /jobs/view/:id
//
//		Restrictions     None
//
// 		Template         base/jobs_view.html
func JobView(ctx echo.Context) error {
	id, err := utils.GetInt(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.BadRequestMessage)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	job, err := query.GetJobByID(id)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	if job != nil {
		utils.SetData(ctx, "Job", job)
	}
	return ctx.Render(http.StatusOK, tmpl.BaseJobsViewTpl, utils.GetData(ctx))
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
	var language string
	switch lang {
	case "en", "sw":
		language = lang
	default:
		language = "en"

	}
	store := session.New()
	sess, _ := store.Get(ctx.Request(), settings.LangSessionName)
	sess.Values[settings.LangSessionKey] = language
	store.Save(ctx.Request(), ctx.Response(), sess)
	ctx.Redirect(http.StatusFound, "/")
	return nil
}

//
//
//		DOCS
//
//

var (
	docIndex  = "docIndex"
	docsRoute = "/docs/"
)

// DocsHome renders the home.md document for the given language.
//
//		Method           GET
//
//		Route            /docs
//
//		Restrictions     None
//
// 		Template         base/docs_index.html
//
func DocsHome(ctx echo.Context) error {
	data := utils.GetData(ctx).(utils.Data)
	lang := data.Get(settings.LangDataKey).(string)
	home := settings.DocsPath + "/" + lang + "/" + settings.DocsIndexPage
	d, err := static.Asset(home)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	data.Set("doc", string(d))
	data.Set(docIndex, getDocIndex(lang))
	data.Set("PageTitle", settings.DocsIndexPage)
	return ctx.Render(http.StatusOK, tmpl.BaseDocsHomeTpl, data)
}

// Docs renders individual zedlist document.
//
//		Method           GET
//
//		Route            /docs/:name
//
//		Restrictions     None
//
// 		Template         base/docs.html
func Docs(ctx echo.Context) error {
	data := utils.GetData(ctx).(utils.Data)
	lang := data.Get(settings.LangDataKey).(string)
	fname := ctx.Param("name")
	if filepath.Ext(fname) != ".md" {
		fname = fname + ".md"
	}
	fPath := settings.DocsPath + "/" + lang + "/" + fname
	d, err := static.Asset(fPath)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	data.Set("doc", string(d))
	data.Set("PageTitle", fname)
	data.Set(docIndex, getDocIndex(lang))
	return ctx.Render(http.StatusOK, tmpl.BaseDocTpl, data)
}

// Doc is a markdown document residing in the static/docs path. It represend documents
// that are internally shipped with zedlist.
type Doc struct {
	Name string
	URL  string
}

func getDocIndex(lang string) []*Doc {
	rst := []*Doc{}
	fdir := settings.DocsPath + "/" + lang
	files, err := static.AssetDir(fdir)
	if err != nil {
		log.Error(nil, err)
		return rst
	}
	for _, v := range files {
		fURL := settings.App.AppURL + docsRoute + v
		name := strings.TrimSuffix(v, filepath.Ext(v))
		rst = append(rst, &Doc{name, fURL})
	}
	return rst
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
	regs, err := query.GetAllRegions()
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
	jobs, count, err := query.GetJobByRegionShort(name)
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
	jobs, err := query.GetJobByRegionPaginate(name, offset, limit)
	if err != nil {
		utils.SetData(ctx, "Message", tmpl.NotFoundMessage)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, settings.JobsListKey, jobs)
	return ctx.Render(http.StatusOK, tmpl.BaseRegionsPaginateTpl, utils.GetData(ctx))
}
