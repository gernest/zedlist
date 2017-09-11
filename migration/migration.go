// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package migration contains migrations for zedlist.
package migration

import (
	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
)

// MigrateTables creates all the database tables used by zedlist.
func MigrateTables() {
	db.Conn.AutoMigrate(&models.Job{})
	db.Conn.AutoMigrate(&models.Region{})
	db.Conn.AutoMigrate(&models.Session{})
	db.Conn.AutoMigrate(&models.User{})
	db.Conn.AutoMigrate(&models.Person{})
	db.Conn.AutoMigrate(&models.PersonName{})
	db.Conn.AutoMigrate(&models.Token{})
	db.Conn.AutoMigrate(&models.Claim{})

	// resume
	db.Conn.AutoMigrate(&models.Resume{})
	db.Conn.AutoMigrate(&models.Award{})
	db.Conn.AutoMigrate(&models.Basic{})
	db.Conn.AutoMigrate(&models.Education{})
	db.Conn.AutoMigrate(&models.Interest{})
	db.Conn.AutoMigrate(&models.Item{})
	db.Conn.AutoMigrate(&models.Language{})
	db.Conn.AutoMigrate(&models.Location{})
	db.Conn.AutoMigrate(&models.SocialProfile{})
	db.Conn.AutoMigrate(&models.Publication{})
	db.Conn.AutoMigrate(&models.Referee{})
	db.Conn.AutoMigrate(&models.Skill{})
	db.Conn.AutoMigrate(&models.Work{})
}

// DropTables drops all the database tables  used by zedlist
func DropTables() {
	db.Conn.DropTableIfExists(&models.Job{})
	db.Conn.DropTableIfExists(&models.Region{})
	db.Conn.DropTableIfExists(&models.Session{})
	db.Conn.DropTableIfExists(&models.User{})
	db.Conn.DropTableIfExists(&models.Person{})
	db.Conn.DropTableIfExists(&models.PersonName{})
	db.Conn.DropTableIfExists(&models.Token{})
	db.Conn.DropTableIfExists(&models.Claim{})
	db.Conn.DropTableIfExists(&models.Resume{})
	db.Conn.DropTableIfExists(&models.Award{})
	db.Conn.DropTableIfExists(&models.Basic{})
	db.Conn.DropTableIfExists(&models.Education{})
	db.Conn.DropTableIfExists(&models.Interest{})
	db.Conn.DropTableIfExists(&models.Item{})
	db.Conn.DropTableIfExists(&models.Language{})
	db.Conn.DropTableIfExists(&models.Location{})
	db.Conn.DropTableIfExists(&models.SocialProfile{})
	db.Conn.DropTableIfExists(&models.Publication{})
	db.Conn.DropTableIfExists(&models.Referee{})
	db.Conn.DropTableIfExists(&models.Skill{})
	db.Conn.DropTableIfExists(&models.Work{})
}
