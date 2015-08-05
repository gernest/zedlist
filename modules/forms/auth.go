// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package forms contains form utilities used by zedlist
package forms

import (
	"time"

	"github.com/gernest/gforms"
	"github.com/gernest/zedlist/modules/i18n"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/melvinmt/gt"
)

var (
	msgRequired  = "message_required"
	msgMinLength = "message_min_length"
	msgEmail     = "message_email"
	msgAge       = "message_age"
	msgEqual     = "message_equal"
)

// Login is the login form
type Login struct {
	Email    string `gforms:"email"`
	Password string `gforms:"password"`
}

// Register is the registration form
type Register struct {
	FirstName       string    `gforms:"first_name"`
	LastName        string    `gforms:"last_name"`
	MiddleName      string    `gforms:"middle_name"`
	Email           string    `gforms:"email"`
	Password        string    `gforms:"password"`
	ConfirmPassword string    `gforms:"confirm_password"`
	Gender          int       `gforms:"gender"`
	BirthDay        time.Time `gforms:"birth_date"`
}

// Form is contains form validation functions, it support translation
// of error messages. This uses gforms.
//
// TODO translate  field names in widgets?
type Form struct {
	tr *gt.Build
}

// New returns a new Form with laguage Target set to lang
func New(lang string) *Form {
	l := i18n.CloneLang()
	l.SetTarget(lang)
	return &Form{l}
}

// LoginForm returns a gform model for Login.
func (f *Form) LoginForm() gforms.ModelForm {
	var attrs = map[string]string{
		"class": "input-large",
	}
	return gforms.DefineModelForm(Login{}, gforms.NewFields(
		gforms.NewTextField(
			"email",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
				gforms.EmailValidator(f.tr.T(msgEmail)),
			},
			gforms.BaseTextWidget("email", attrs),
		),
		gforms.NewTextField(
			"password",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
				gforms.MinLengthValidator(6, f.tr.T(msgMinLength, 6)),
			},
			gforms.PasswordInputWidget(attrs),
		),
	))
}

// RegisterForm implements gforms.ModelForm interface for Registration form.
func (f *Form) RegisterForm() gforms.ModelForm {
	var birtdateAttrs = map[string]string{
		"id":    "birthdate",
		"class": "input-large",
	}
	var inputAttrs = map[string]string{
		"class": "input-large",
	}
	return gforms.DefineModelForm(Register{}, gforms.NewFields(
		gforms.NewTextField(
			"first_name",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.TextInputWidget(inputAttrs),
		),
		gforms.NewTextField(
			"last_name",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.TextInputWidget(inputAttrs),
		),
		gforms.NewTextField(
			"middle_name",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.TextInputWidget(inputAttrs),
		),
		gforms.NewTextField(
			"email",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
				gforms.EmailValidator(f.tr.T(msgEmail)),
			},
			gforms.BaseTextWidget("email", inputAttrs),
		),
		gforms.NewTextField(
			"password",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
				gforms.MinLengthValidator(6, f.tr.T(msgMinLength, 6)),
			},
			gforms.PasswordInputWidget(inputAttrs),
		),
		gforms.NewTextField(
			"confirm_password",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
				gforms.MinLengthValidator(6, f.tr.T(msgMinLength, 6)),
				EqualValidator{To: "password", Message: f.tr.T(msgEqual)},
			},
			gforms.PasswordInputWidget(inputAttrs),
		),
		gforms.NewTextField(
			"gender",
			gforms.Validators{
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.SelectWidget(
				inputAttrs,
				func() gforms.SelectOptions {
					return gforms.StringSelectOptions([][]string{
						{"Select...", "", "true", "false"},
						{"Male", "0", "false", "false"},
						{"Female", "1", "false", "false"},
						{"Zombie", "2", "false", "true"},
					})
				},
			),
		),

		gforms.NewDateTimeField(
			"birth_date",
			settings.App.BirthDateFormat,
			gforms.Validators{
				BirthDateValidator{Limit: settings.App.MinimumAge, Message: f.tr.T(msgAge, settings.App.MinimumAge)},
				gforms.Required(f.tr.T(msgRequired)),
			},
			gforms.BaseTextWidget("text", birtdateAttrs),
		),
	))
}
