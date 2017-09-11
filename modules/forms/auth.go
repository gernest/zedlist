// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package forms contains form utilities used by zedlist
package forms

import (
	"net/http"
	"strings"

	"github.com/gernest/gt"
	"github.com/gorilla/schema"
	"github.com/zedio/zedlist/modules/i18n"
)

var (
	msgRequired  = "message_required"
	msgMinLength = "message_min_length"
	msgEmail     = "message_email"
	msgAge       = "message_age"
	msgEqual     = "message_equal"
)

var decoder = schema.NewDecoder()

// Login is the login form
type Login struct {
	Name     string `schema:"username"`
	Password string `schema:"password"`
	CSRF     string `schema:"csrf_token"`
	vals     map[string]string
}

func (l *Login) Valid() bool {
	if l.vals == nil {
		l.vals = make(map[string]string)
	}
	var hasError bool
	l.vals["username"] = l.Name
	l.vals["password"] = l.Password
	l.Name = strings.TrimSpace(l.Name)
	if l.Name == "" {
		l.vals["username_error"] = "username is required"
		hasError = true
	}
	if l.Password == "" {
		l.vals["password_error"] = "password is required"
		hasError = true
	}
	return !hasError
}
func (l *Login) Ctx() map[string]string {
	return l.vals
}

// Register is the registration form
type Register struct {
	UserName        string `schema:"username"`
	Email           string `schema:"email"`
	Password        string `schema:"password"`
	ConfirmPassword string `schema:"confirm_password"`
	CSRF            string `schema:"csrf_token"`
	vals            map[string]string
}

func (r *Register) Valid() bool {
	if r.vals == nil {
		r.vals = make(map[string]string)
	}
	var haserror bool
	r.vals["username"] = r.UserName
	r.vals["email"] = r.Email
	r.vals["password"] = r.Password
	r.vals["confirm_password"] = r.ConfirmPassword
	r.UserName = strings.TrimSpace(r.UserName)
	if r.UserName == "" {
		r.vals["username_error"] = "username is required"
		haserror = true
	}
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		r.vals["email_error"] = "email is required"
		haserror = true
	}
	r.Password = strings.TrimSpace(r.Password)
	if r.Password == "" {
		r.vals["password_error"] = "password is required"
		haserror = true
	}
	r.ConfirmPassword = strings.TrimSpace(r.ConfirmPassword)
	if r.ConfirmPassword == "" {
		r.vals["confirm_password_error"] = "confirm password is required"
		haserror = true
	} else {
		if r.Password != r.ConfirmPassword {
			r.vals["confirm_password_error"] = "confirm password must match passwordd"
			haserror = true
		}
	}
	return !haserror
}

func (r *Register) Ctx() map[string]string {
	return r.vals
}

// Form is contains form validation functions, it support translation
// of error messages. This uses schema.
//
// TODO translate  field names in widgets?
type Form struct {
	tr *gt.Build
	l  *Login
	r  *Register
}

// New returns a new Form with laguage Target set to lang
func New(lang string) *Form {
	l := i18n.CloneLang()
	l.SetTarget(lang)
	return &Form{tr: l}
}

func (f *Form) SetLogin(l *Login) {
	f.l = l
}

func (f *Form) DecodeLogin(r *http.Request) (*Login, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	l := &Login{}
	if err := decoder.Decode(l, r.PostForm); err != nil {
		return nil, err
	}
	return l, nil
}

func (f *Form) SetRegister(r *Register) {
	f.r = r
}

func (f *Form) DecodeRegister(r *http.Request) (*Register, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	l := &Register{}
	if err := decoder.Decode(l, r.PostForm); err != nil {
		return nil, err
	}
	return l, nil
}

func (f *Form) DecodeDelete(r *http.Request) string {
	r.ParseForm()
	return r.Form.Get("delete")
}
