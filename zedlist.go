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
	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/gernest/zedlist/modules/tmpl"

	"github.com/gernest/zedlist/routes/auth"
	"github.com/gernest/zedlist/routes/base"
	"github.com/gernest/zedlist/routes/dash"
	"github.com/gernest/zedlist/routes/japi"
	"github.com/gernest/zedlist/routes/resume"
	"github.com/gernest/zedlist/routes/search"

	"github.com/gernest/zedlist/migration"

	"github.com/codegangsta/cli"
	ass "github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Routes is the echo instance with all zedlist routes.
func Routes() *echo.Echo {
	e := echo.New()

	e.SetRenderer(tmpl.NewRenderer()) // Set renderer

	// middlewares
	e.Use(i18n.Langs())      // languages
	e.Use(flash.Flash())     // flash messages
	e.Use(userAuth.Normal()) // adding user context data
	e.Use(middleware.Gzip()) // Gzip

	// API
	a := e.Group("/api")
	a.Post("/jobs", japi.CreateJob)
	a.Get("/jobs/:id", japi.GetJob)
	a.Get("/jobs", japi.GetIndex)
	a.Put("/jobs", japi.UpdateJob)

	// HOME
	e.Get("/", base.Home)
	e.Get("/language/:lang", base.SetLanguage)

	// DOCS
	e.Get("/docs", base.DocsHome)
	e.Get("/docs/:name", base.Docs)

	// BASE
	b := e.Group("/jobs")
	b.Get("/", base.JobsHome)
	b.Get("/view/:id", base.JobView)
	b.Get("/regions", base.RegionsHome)
	b.Get("/regions/:name", base.RegionsJobView)
	b.Get("/regions/:name/:from/:to", base.RegionsJobPaginate)

	// AUTH
	xauth := e.Group("/auth")

	// add csrf protection
	xauth.Use(csrf.Nosurf())
	xauth.Use(csrf.Tokens())

	xauth.Get("/login", auth.Login)
	xauth.Post("/login", auth.LoginPost)
	xauth.Get("/register", auth.Register)
	xauth.Post("/register", auth.RegisterPost)
	xauth.Get("/logout", auth.Logout)

	// DASHBOARD
	dashboard := e.Group("/dash")
	dashboard.Use(userAuth.Must())
	dashboard.Get("/", dash.Home)
	dashboard.Get("/jobs/new", dash.JobsNewGet)
	dashboard.Post("/jobs/new", dash.JobsNewPost)
	dashboard.Get("/profile", dash.Profile)
	dashboard.Post("/profile/name", dash.ProfileName)

	// RESUME
	r := e.Group("/resume")
	r.Get("/", resume.Home)
	r.Post("/new", resume.Create)
	r.Get("/view:id", resume.View)
	r.Post("/update/:id", resume.Update)
	r.Post("/delete/:id", resume.Delete)

	// SEARCH
	e.Post("/search", search.Find)

	// STATIC
	box := &ass.AssetFS{
		Asset:    static.Asset,
		AssetDir: static.AssetDir,
		Prefix:   "static",
	}
	staticFileServer := http.StripPrefix("/static/", http.FileServer(box))
	e.Get("/static/*", staticFileServer)
	return e
}

// Server runs the zedlist server.
func Server(ctx *cli.Context) {
	log.Info(nil, fmt.Sprintf("starting zedlist server at %s", settings.App.AppURL))
	r := Routes()
	r.Run(fmt.Sprintf(":%d", settings.App.Port))
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
	query.PopulateDB()
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
