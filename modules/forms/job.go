// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import "net/http"

// JobForm is a form for job posting
type JobForm struct {
	Title       string `schema:"title"`
	Description string `schema:"description"`
	CSRF        string `schema:"csrf_token"`
}

func (j *JobForm) Valid() bool {
	return true
}

// JobForm implements gform.ModelForm interface.
func (f *Form) DecodeJobForm(r *http.Request) (*JobForm, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	l := &JobForm{}
	if err := decoder.Decode(l, r.PostForm); err != nil {
		return nil, err
	}
	return l, nil
}
