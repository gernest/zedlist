// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	"testing"

	"github.com/gernest/zedlist/models"
)

func TestValidatePersonNames(t *testing.T) {
	p := &models.PersonName{
		FamilyName: "ernest",
		GivenName:  "geofrey",
		MiddleName: "katesigwa",
	}
	v := NewValid("en")
	errs := v.ValidatePersonName(p)
	if errs != nil {
		t.Errorf("expected nil got %v", errs)
	}

	p.GivenName = ""
	errs = v.ValidatePersonName(p)
	if errs != nil {
		t.Errorf("expected nil got %v", errs)
	}

	// Invalid
	p.FamilyName = "---@**"
	errs = v.ValidatePersonName(p)
	if errs == nil {
		t.Error("expected errors got nil")
	}
}
