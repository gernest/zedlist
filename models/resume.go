// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"time"
)

// Resume represent Curriculum Vitae
type Resume struct {
	ID           int64         `json:"id"`
	PersonID     int64         `json:"person_id"`
	Title        string        `json:"title"`
	Basic        Basic         `json:"basic"`
	BasicID      int64         `json:"basic_id"`
	Work         []Work        `json:"work"`
	Volunteer    []Work        `json:"volunteer"`
	Education    []Education   `json:"education"`
	Awards       []Award       `json:"awards"`
	Publications []Publication `json:"publications"`
	Skills       []Skill       `json:"skills"`
	Languages    []Language    `json:"languages"`
	Interests    []Interest    `json:"interests"`
	References   []Referee     `json:"references"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// Basic is the basic information for a resume
type Basic struct {
	ID int64 `json:"id" schema:"-"`

	// Name is the name of the resume holder
	// e.g Geofrey Ernest
	Name string `json:"name" schema:"name"`

	// Label is a short description
	// e.g A dreamer from Tanzania
	Label string `json:"label" schema:"label"`

	// Picture is a URL to an image file, format JPEG or PNG
	Picture string `json:"picture" schema:"picture"`

	// Email is an email address
	// e.g gernest@wakalikwanza.tz
	Email string `json:"email" schema:"email"`

	// Phone is the phonenumber, this is used as a string
	// to suppoert mutliple formats.
	Phone string `json:"phone" schema:"phone"`

	// Website is the URL to a website.
	Website string `json:"website" schema:"website"`

	//Summary is a 2-3 sentence biography about resume owner
	Summary string `json:"summary" schema:"summary"`

	ResumeLocation   Location `json:"location" schema:"location"`
	ResumeLocationID int64    `json:"-" shema:"-"`

	Profiles []SocialProfile `json:"profiles" schema:"-"`

	CreatedAt time.Time `json:"created_at" schema:"-"`
	UpdatedAt time.Time `json:"updated_at" schema:"-"`
}

// Location is the location details of a resume owner.
type Location struct {
	ID          int64     `json:"id" schema:"-"`
	Address     string    `json:"address" schema:"address"`
	PostalCode  string    `json:"postalCode" schema:"postal_code"`
	City        string    `json:"city" schema:"city"`
	CountryCode string    `json:"countryCode" schema:"country_code"`
	Region      string    `json:"region" schema:"region"`
	CreatedAt   time.Time `json:"created_at" schema:"-"`
	UpdatedAt   time.Time `json:"updated_at" schema:"-"`
}

// SocialProfile is the profile of the resume owner.
type SocialProfile struct {
	ID            int64 `json:"id"`
	ResumeBasicID int64 `sql:"index"`

	// Network is a social network
	// e.g Facebook google+ or github
	Network string `json:"network"`

	// UserName is the social network's username
	// e.g gerest
	UserName string `json:"username"`

	// URL is the link to the profile
	// e.g https://github.com/gernest
	URL string `json:"string"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Work is the work details of the resume owner.
type Work struct {
	ID       int64 `json:"id"`
	ResumeID int64 `sql:"index"`

	// e.g zedlist
	Company string `json:"company"`

	// e.g founder
	Position string `json:"position"`

	//e.g zedlist.co.tz
	Website string `json:"website"`

	// StartDate is the time you started working in the company
	// TODO(gernest): use a string?
	StartDate time.Time `json:"startDate"`

	// EndDate is the time you stoped working for the company
	// TODO(gernest): use string?
	EndDate time.Time `json:"endDate"`

	// Summary is the overview of your responsibilities at the company.
	Summary string `json:"summary"`

	//Highlights is a list of accomplishments
	Highlights []Item `json:"highlights" gorm:"polymorphic:Item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Item is a single item.
type Item struct {
	ID        int64     `json:"id"`
	Body      string    `json:"item"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Education is the education details of the resume owner.
type Education struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`

	// e.g UDOM
	Institution string `json:"institution"`

	// e.g Business
	Area string `json:"area"`

	// e.g Bachelor
	StudyType string `json:"stydyType"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	// eg 4.0
	GPA string `json:"gpa"`

	// e.g Development studies
	Courses []Item `json:"courses" gorm:"polymorphic:Item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Award is the award details of the resume owner.
type Award struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`
	// e.g Best dreamer
	Title string `json:"title"`

	// e.g 2015-10-10
	Date time.Time `json:"date"`

	// e.g  my uncle
	Awarder string `json:"awarder"`

	// e.g  A big time dreamer in a big dream
	Summary string `json:"summary"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Publication is the publication details of the resume owner.
type Publication struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`

	// e.g Beating poverty with golang
	Name string `json:"name"`

	// e.g Zedlist
	Publisher string `json:"publisher"`

	// e.g 2005-10-10
	ReleaseDate time.Time `json:"releaseDate"`

	// e.g zedlist.com
	Website string `json:"website"`

	// e.g cheating on poverty with code
	Summary string `json:"summary"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Skill is the skill details of the resume owner.
type Skill struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`

	// e.g web dev
	Name string `json:"name"`

	// e.g Intermediary
	Level string `json:"level"`

	// e.g Golang
	KeyWords []Item `json:"keywords" gorm:"polymorphic:Item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Language is the language details of the resume owner.
type Language struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`

	// e.g English, Swahili
	Language string `json:"language"`

	// e.g fluent, beginner
	Fluency string `json:"fluency"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Interest is the interest details of the resume owner.
type Interest struct {
	ID       int64 `json:"id"`
	ResumeID int64 `json:"resume_id"`

	// e.g music
	Name string `json:"name"`

	// e.g Vampire Weekend
	Keywords []Item `json:"keywords" gorm:"polymorphic:Item"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Referee is the details about a person who is refereeing for the resume owner.
type Referee struct {
	ID        int64     `json:"id"`
	ResumeID  int64     `json:"resume_id"`
	Name      string    `json:"name"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//SampleResume cretes a sample resume.
func SampleResume() Resume {
	now := time.Now()
	return Resume{
		Title: "Sample Resume",
		Basic: Basic{
			Name:    "John Doe",
			Label:   "A dreamer from tz",
			Picture: "/static/img/gopher.png",
			Email:   "johndoe@zedlist.io",
			Phone:   "+2557690000000",
			Website: "www.zedlist.com",
			Summary: "WITNESS NE",
			ResumeLocation: Location{
				Address:     "home",
				PostalCode:  "200000",
				City:        "Mwanza",
				CountryCode: "+255",
				Region:      "Mwanza",
			},
			Profiles: []SocialProfile{
				SocialProfile{
					Network:  "github",
					UserName: "johnDoe",
					URL:      "gihub.com/johnDoe",
				},
			},
		},
		Work: []Work{
			Work{
				Company:   "johnDee",
				Position:  "founder",
				Website:   "johndoe.com",
				StartDate: now,
				EndDate:   now,
				Summary:   "WITNESS ME",
				Highlights: []Item{
					Item{
						Body: "paid to sleep",
					},
				},
			},
		},
		Volunteer: []Work{
			Work{
				Company:   "johnDee",
				Position:  "founder",
				Website:   "johndoe.com",
				StartDate: now,
				EndDate:   now,
				Summary:   "WITNESS ME",
				Highlights: []Item{
					Item{
						Body: "charity and stuffs",
					},
				},
			},
		},
		Education: []Education{
			Education{
				Institution: "UFOS",
				Area:        "busness",
				StudyType:   "bachelor",
				GPA:         "4.0",
				StartDate:   now,
				EndDate:     now,
				Courses: []Item{
					Item{
						Body: "whacko",
					},
				},
			},
		},
		Awards: []Award{
			Award{
				Title:   "dreamer",
				Date:    now,
				Awarder: "Van Helsing",
				Summary: "best dreamer",
			},
		},
		Publications: []Publication{
			Publication{
				Name:        "Witness Me",
				Publisher:   "whacko",
				ReleaseDate: now,
				Website:     "whacko.io",
				Summary:     "WITNESS ME",
			},
		},
		Skills: []Skill{
			Skill{
				Name:  "web -dev",
				Level: "master",
				KeyWords: []Item{
					Item{
						Body: "golang",
					},
				},
			},
		},
		Languages: []Language{
			Language{
				Language: "swahili",
				Fluency:  "master",
			},
		},
		Interests: []Interest{
			Interest{
				Name: "music",
				Keywords: []Item{
					Item{
						Body: "vampire weekend",
					},
				},
			},
		},
		References: []Referee{
			Referee{
				Name:      "whacko wacka",
				Reference: " grave digger",
			},
		},
	}
}
