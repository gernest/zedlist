// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	validate "github.com/asaskevich/govalidator"
	"github.com/gernest/zedlist/models"

	"github.com/gernest/gt"
	"github.com/gernest/vala"
	"github.com/gernest/zedlist/modules/i18n"
)

var (
	validNameMsg = "valid_name_msg"
)

// Valid is a validation struct for form params
type Valid struct {
	lang *gt.Build
}

//NewValid creates a new Valid instance
func NewValid(lang string) *Valid {
	l := i18n.CloneLang()
	l.SetTarget(lang)
	return &Valid{l}
}

// Name checks names
func (v *Valid) Name(name, field string) vala.Checker {
	return func() (bool, string) {
		return isName(name), v.lang.T(validNameMsg, v.lang.T(field))
	}
}

func isName(name string) bool {
	if name != "" {
		return validate.IsAlpha(name)
	}
	return true
}

// ValidatePersonName validates person name and returns a slice of error messages
// to be shown to the user.
func (v *Valid) ValidatePersonName(p *models.PersonName) []string {
	return checkErr(vala.BeginValidation().Validate(
		v.Name(p.GivenName, "first_name"),
		v.Name(p.FamilyName, "family_name"),
		v.Name(p.MiddleName, "middle_	name"),
	))
}

func checkErr(v *vala.Validation) []string {
	if v == nil || len(v.Errors) <= 0 {
		return nil
	}
	return v.Errors
}
