// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models contains models for zedlist
package models

import (
	"fmt"
	"time"

	"github.com/gernest/zedlist/modules/sanitize"
)

// JSONError is json resopnse error
type JSONError struct {
	Message string `json:"error"`
}

// NewJSONErr creates a new JSONERRor
func NewJSONErr(msg string) *JSONError {
	return &JSONError{msg}
}

// Job is the base job struct
type Job struct {
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title"`

	// Description is a longer description about the job
	// the text format should be in Markdown.
	Description string `json:"description" sql:"null;type:text"`

	PersonID int64

	Region    Region    `json:"region"`
	RegionID  int64     `json:"region_id"`
	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Sanitize we don't trust the user, so we clean up.
func (j *Job) Sanitize() {
	// use in house sanitization policy for the title.
	j.Title = sanitize.Title(j.Title)
}

// Logger is the logging infomation schema.
type Logger struct {
	ID        int64
	Level     int
	Message   string
	Path      string
	User      string
	Time      time.Time
	CreatedAt time.Time
	UpdateAt  time.Time
}

const (

	// ObjPerson is the user of type person.
	ObjPerson = "person"

	// ObjOrganization is the user of type organization
	ObjOrganization = "organization"

	//Single is for not in any relationships
	Single = "single"

	// Married is married
	Married = "married"

	// Complicated is in relation which is complicated
	Complicated = "complicated"

	// OpenRelationship is in open srelationship
	OpenRelationship = "open"

	// Widowed is in a situation with a decesed spouse
	Widowed = "widowed"

	// DomesticRelation is in domestic relationship
	DomesticRelation = "domestic"

	// CivilUnion is in civil union
	CivilUnion = "civil"

	male = iota
	female
	zombie
)

// Person contains user's details
type Person struct {
	ID                 int64
	AboutMe            string
	Birthday           time.Time
	CurrentLocation    string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DisplayName        string
	Email              string
	Jobs               []Job
	Gender             int
	PersonName         PersonName
	PersonNameID       int
	ObjectType         string
	RelationshipStatus string
	URL                string
	Resumes            []Resume
}

// PersonName contains user's names.
type PersonName struct {
	ID         int64  `shcema:"-"`
	Name       string `shcema:"name"`
	FamilyName string `schema:"family_name"`
	GivenName  string `schema:"given_name"`
	MiddleName string `schema:"middle_name"`
}

// Formatted returns the person's full name as a sigle string
func (p *PersonName) Formatted() string {
	return fmt.Sprintf("%s %s %s", p.GivenName, p.MiddleName, p.FamilyName)
}

// UpdateNames updates profile names.
func (p *Person) UpdateNames(pName *PersonName) {
	if pName.FamilyName != "" {
		p.PersonName.FamilyName = sanitize.Name(pName.FamilyName)
	}
	if pName.MiddleName != "" {
		p.PersonName.MiddleName = sanitize.Name(pName.MiddleName)
	}
	if pName.GivenName != "" {
		p.PersonName.GivenName = sanitize.Name(pName.GivenName)
	}
}

// Region represent the region(look I'm from Tanzania, here we have cities and regions).
type Region struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

// Session stores session data from gorilla/sessions
type Session struct {
	ID        int64
	Key       string
	Data      string `sql:"type:text"`
	ExpiresOn time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Token represent a jwt token object.
type Token struct {
	ID        int64
	Key       string `sql:"unique_index;type:text"`
	Valid     bool
	Claims    []Claim
	ExpireOn  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Claim represent claims to a token with a given TokenID
type Claim struct {
	ID      int64
	TokenID int `sql:"index"`
	Key     string
	Value   string
}

const (

	// StatusActive indicates that the user is free to interact with the app.
	StatusActive = iota

	// StatusSuspended indicates a user has been restricted to interact with the app
	// temporarilly.
	StatusSuspended

	// StatusDeleted indicate s the user is no longer a member of the app.
	StatusDeleted
)

// User is the zedlist user.
type User struct {
	ID        int64
	Email     string `sql:"type:varchar(100);unique_index"`
	Password  string
	Name      string `gorm:"unique"`
	Person    Person
	PersonID  int
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
