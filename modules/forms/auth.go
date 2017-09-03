// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package forms contains form utilities used by zedlist
package forms

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gernest/gt"
	"github.com/gernest/zedlist/modules/i18n"
	"github.com/gorilla/schema"
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
}

func (l *Login) Valid() bool {
	return true
}

// Register is the registration form
type Register struct {
	UserName        string `schema:"username"`
	Email           string `schema:"email"`
	Password        string `schema:"password"`
	ConfirmPassword string `schema:"confirm_password"`
	CSRF            string `schema:"csrf_token"`
}

func (r *Register) Valid() bool {
	return true
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

func (f *Form) Login() template.HTML {
	return template.HTML(fmt.Sprintf(`
	<div class="field">
		<label> %s </label>
		<input type="text" name="username" placeholder="%s">
	</div>
	<div class="field">
		<label> %s</label>
		<input type="password" name="password" placeholder="%s">
	</div>
	<button class="ui fluid large submit button"> %s</button>
`, f.tr.T("username"), f.tr.T("username_or_email"),
		f.tr.T("password"), f.tr.T("password"),
		f.tr.T("login"),
	))
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

func (f *Form) Register() template.HTML {
	return template.HTML(fmt.Sprintf(`
	<div class="field">
		<label> %s </label>
		<input type="text" name="username" placeholder="%s">
	</div>
	<div class="field">
		<label> %s </label>
		<input type="text" name="email" placeholder="%s">
	</div>
	<div class="field">
		<label> %s</label>
		<input type="password" name="password" placeholder="%s">
	</div>
	<div class="field">
		<label> %s</label>
		<input type="password" name="confirm_password" placeholder="%s">
	</div>
	<button class="ui fluid large submit button"> %s</button>
`, f.tr.T("username"), f.tr.T("username"),
		f.tr.T("email"), f.tr.T("email"),
		f.tr.T("password"), f.tr.T("password"),
		f.tr.T("confirm_password"), f.tr.T("confirm_password"),
		f.tr.T("register"),
	))
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

func (f *Form) Delete() template.HTML {
	return template.HTML(`
<div class="field">
	<input type="text" name="delete" placeholder="Type your username">
</div>	
	`)
}

func (f *Form) DecodeDelete(r *http.Request) string {
	r.ParseForm()
	return r.Form.Get("delete")
}
