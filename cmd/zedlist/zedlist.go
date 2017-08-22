// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/gernest/zedlist/modules/version"

	app "github.com/gernest/zedlist"
	"github.com/urfave/cli"
)

func main() {
	zedApp := cli.NewApp()
	zedApp.Name = "Zedlist"
	zedApp.Usage = "A job recruitment service"
	zedApp.Version = version.VERSION
	zedApp.Authors = app.Authors
	zedApp.Commands = []cli.Command{
		app.ServerCommand,
		app.MigrateCommand,
	}
	zedApp.Run(os.Args)
}
