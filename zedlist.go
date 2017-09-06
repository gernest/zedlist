// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package zedlist is a job recruitment service.
package zedlist

import (
	"fmt"
	"net/http"

	userAuth "github.com/gernest/zedlist/middlewares/auth"
	"github.com/gernest/zedlist/middlewares/csrf"
	"github.com/gernest/zedlist/middlewares/flash"
	"github.com/gernest/zedlist/middlewares/i18n"

	"github.com/gernest/zedlist/bindata/static"
	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"

	"github.com/gernest/zedlist/routes/auth"
	"github.com/gernest/zedlist/routes/base"
	"github.com/gernest/zedlist/routes/japi"
	"github.com/gernest/zedlist/routes/jobs"
	"github.com/gernest/zedlist/routes/resume"
	"github.com/gernest/zedlist/routes/search"

	"github.com/gernest/zedlist/migration"

	ass "github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

// Routes is the echo instance with all zedlist routes.
func Routes() *echo.Echo {
	e := echo.New()

	e.Renderer = tmpl.NewRenderer() // Set renderer

	// middlewares
	e.Use(utils.WrapMiddleware(i18n.Langs()))      // languages
	e.Use(utils.WrapMiddleware(flash.Flash()))     // flash messages
	e.Use(utils.WrapMiddleware(userAuth.Normal())) // adding user context data
	e.Use(middleware.Gzip())                       // Gzip

	// API
	a := e.Group("/api")
	a.POST("/jobs", japi.CreateJob)
	a.GET("/jobs/:id", japi.GetJob)
	a.GET("/jobs", japi.GetIndex)
	a.PUT("/jobs", japi.UpdateJob)

	// HOME
	e.GET("/", base.Home)
	e.GET("/language/:lang", base.SetLanguage)

	// DOCS
	e.GET("/docs", base.DocsHome)
	e.GET("/docs/:name", base.Docs)

	// AUTH
	xauth := e.Group("/auth")

	// add csrf protection
	xauth.Use(echo.WrapMiddleware(csrf.Nosurf()))
	xauth.Use(utils.WrapMiddleware(csrf.Tokens()))

	xauth.GET("/login", auth.Login)
	xauth.POST("/login", auth.LoginPost)
	xauth.GET("/register", auth.Register)
	xauth.POST("/register", auth.RegisterPost)
	xauth.GET("/logout", auth.Logout)
	xauth.GET("/delete", auth.Delete)
	xauth.POST("/delete", auth.DeletePost)

	// JOBS
	j := e.Group("/jobs")
	j.Use(utils.WrapMiddleware(userAuth.Must()))
	j.GET("/", jobs.List)
	j.GET("/new", jobs.New)
	j.POST("/new", jobs.NewPost)
	j.GET("/view/:id", jobs.View)

	// RESUME
	r := e.Group("/resume")
	r.GET("/", resume.Home)
	r.POST("/new", resume.Create)
	r.GET("/view:id", resume.View)
	r.POST("/update/:id", resume.Update)
	r.POST("/delete/:id", resume.Delete)

	// SEARCH
	e.POST("/search", search.Find)

	// STATIC
	box := &ass.AssetFS{
		Asset:    static.Asset,
		AssetDir: static.AssetDir,
		Prefix:   "static",
	}
	staticFileServer := http.StripPrefix("/static/", http.FileServer(box))
	e.GET("/static/*", echo.WrapHandler(staticFileServer))
	return e
}

// Server runs the zedlist server.
func Server(ctx *cli.Context) {
	log.Info(nil, fmt.Sprintf("starting zedlist server at %s", settings.App.AppURL))
	r := Routes()
	r.Start(fmt.Sprintf(":%d", settings.App.Port))
}

// Authors are the authors of zedlist
var Authors = []cli.Author{
	{"Geofrey Ernest", "geofreyernest@live.com"},
}

// ServerCommand for running zedlist server
var ServerCommand = cli.Command{
	Name:        "server",
	ShortName:   "s",
	Usage:       "Runs zedlist server",
	Description: `starts a loal web server`,
	Action:      Server,
}

// Migrate runs migrations
func Migrate(ctx *cli.Context) {
	if ctx.BoolT("d") {
		// run migrations
		migration.DropTables()
	}
	migration.MigrateTables()
	query.PopulateDB(db.Conn)
}

// FlagDev is true when in developmet mode.
var FlagDev = cli.BoolTFlag{
	Name:   "d",
	Usage:  "true if in development mode",
	EnvVar: "DEV_MODE",
}

// MigrateCommand is a command for running migrations.
var MigrateCommand = cli.Command{
	Name:        "migrate",
	ShortName:   "m",
	Usage:       "Runs migrations for zedlist",
	Description: `creates database tables and populate them with data if necessary`,
	Action:      Migrate,
	Flags: []cli.Flag{
		FlagDev,
	},
}
