// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
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
	db.Conn.AutoMigrate(&models.ResumeAward{})
	db.Conn.AutoMigrate(&models.ResumeBasic{})
	db.Conn.AutoMigrate(&models.ResumeEducation{})
	db.Conn.AutoMigrate(&models.ResumeInterest{})
	db.Conn.AutoMigrate(&models.ResumeItem{})
	db.Conn.AutoMigrate(&models.ResumeLanguage{})
	db.Conn.AutoMigrate(&models.ResumeLocation{})
	db.Conn.AutoMigrate(&models.ResumeProfile{})
	db.Conn.AutoMigrate(&models.ResumePublication{})
	db.Conn.AutoMigrate(&models.ResumeReferee{})
	db.Conn.AutoMigrate(&models.ResumeSkill{})
	db.Conn.AutoMigrate(&models.ResumeWork{})
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
	db.Conn.DropTableIfExists(&models.ResumeAward{})
	db.Conn.DropTableIfExists(&models.ResumeBasic{})
	db.Conn.DropTableIfExists(&models.ResumeEducation{})
	db.Conn.DropTableIfExists(&models.ResumeInterest{})
	db.Conn.DropTableIfExists(&models.ResumeItem{})
	db.Conn.DropTableIfExists(&models.ResumeLanguage{})
	db.Conn.DropTableIfExists(&models.ResumeLocation{})
	db.Conn.DropTableIfExists(&models.ResumeProfile{})
	db.Conn.DropTableIfExists(&models.ResumePublication{})
	db.Conn.DropTableIfExists(&models.ResumeReferee{})
	db.Conn.DropTableIfExists(&models.ResumeSkill{})
	db.Conn.DropTableIfExists(&models.ResumeWork{})
}
