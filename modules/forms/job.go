// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	"github.com/gernest/gforms"
)

// JobForm is a form for job posting
type JobForm struct {
	Title       string `gforms:"title"`
	Description string `gforms:"description"`
}

// JobForm implements gform.ModelForm interface.
func (f *Form) JobForm() gforms.ModelForm {
	titleAttrs := map[string]string{
		"id": "job-title",
	}
	descAttr := map[string]string{
		"id": "job-description",
	}
	return gforms.DefineModelForm(JobForm{}, gforms.NewFields(
		gforms.NewTextField(
			"title",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.TextInputWidget(titleAttrs),
		),
		gforms.NewTextField(
			"description",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.TextAreaWidget(descAttr),
		),
	))
}
