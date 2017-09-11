// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package db exposes database connections for zedlist.
package db

import (
	"github.com/jinzhu/gorm"
	"github.com/zedio/zedlist/modules/log"
	"github.com/zedio/zedlist/modules/settings"

	// use postgres only
	_ "github.com/lib/pq"
)

//BaseDB is the base database instance used by zedlist
var BaseDB *DB

// Conn is the global database connection.
var Conn *gorm.DB

func init() {
	BaseDB = NewDB(settings.App)
	d, err := BaseDB.Connect()
	if err != nil {
		log.Error(nil, err)
	}
	Conn = d
}

// DB contains information for accessing a database
type DB struct {
	Diealect string
	ConnStr  string
}

// NewDB creates a new DB instance.
func NewDB(cfg *settings.Config) *DB {
	return &DB{cfg.DBDialect, cfg.DBConn}
}

// Connect establish a database connection and returns a *gorm.DB.
func (db *DB) Connect() (*gorm.DB, error) {
	gdb, err := gorm.Open(db.Diealect, db.ConnStr)
	if err != nil {
		return nil, err
	}
	return gdb, nil
}
