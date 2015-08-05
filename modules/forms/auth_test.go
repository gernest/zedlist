// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestLoginForm(t *testing.T) {
	f := New("sw")

	// Case correct details
	vars := url.Values{
		"email":    {"gernest@mwnza.com"},
		"password": {"mypassword"},
	}
	req1, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form1 := f.LoginForm()(req1)
	if !form1.IsValid() {
		t.Errorf("validating form %v", form1.Errors())
	}
	l := form1.GetModel().(Login)
	if l.Email != vars.Get("email") {
		t.Errorf("retrieving model form: expecting %s got %s", vars.Get("email_address"), l.Email)
	}

	// Case wrong email
	vars.Set("email", "kilimanjaro")
	req2, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form2 := f.LoginForm()(req2)
	if form2.IsValid() {
		t.Error("expected validation error")
	}
	errs := form2.Errors()
	emailErrs := errs.Get("email")
	if !strings.Contains(emailErrs[0], f.tr.T(msgEmail)) {
		t.Errorf("expected %v to contain %v", errs, f.tr.T(msgEmail))
	}

	// Case wrong short password
	vars.Set("email", "kilimanjaro@example.com")
	vars.Set("password", "short")
	req3, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req3.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form3 := f.LoginForm()(req3)
	if form3.IsValid() {
		t.Error("expected validation error")
	}
	errs = form3.Errors()
	passErrs := errs.Get("password")
	if !strings.Contains(passErrs[0], f.tr.T(msgMinLength, 6)) {
		t.Errorf("expected %v to contain %v", errs, f.tr.T(msgMinLength))
	}
}
func TestRendering(t *testing.T) {
	actualRegTpl := "testdata/actual/register.html"
	expectedRegTpl := "testdata/expect/register.html"
	tpl, err := template.ParseFiles(actualRegTpl)
	if err != nil {
		t.Error(err)
	}
	out := &bytes.Buffer{}
	fm := New("en")
	err = tpl.Execute(out, fm.RegisterForm()())
	if err != nil {
		t.Error(err)
	}
	//ioutil.WriteFile(expectedRegTpl, out.Bytes(), 0600)
	expect, err := ioutil.ReadFile(expectedRegTpl)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expect, out.Bytes()) {
		t.Errorf("expected %s \n got %s \n", expect, out.String())
	}

}
func TestRenderingJobForm(t *testing.T) {
	actualRegTpl := "testdata/actual/job.html"
	expectedRegTpl := "testdata/expect/job.html"
	tpl, err := template.ParseFiles(actualRegTpl)
	if err != nil {
		t.Error(err)
	}
	out := &bytes.Buffer{}
	fm := New("en")
	err = tpl.Execute(out, fm.JobForm()())
	if err != nil {
		t.Error(err)
	}
	//ioutil.WriteFile(expectedRegTpl, out.Bytes(), 0600)
	expect, err := ioutil.ReadFile(expectedRegTpl)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expect, out.Bytes()) {
		t.Errorf("expected %s \n got %s \n", expect, out.String())
	}
}
func TestRegisterForm(t *testing.T) {
	f := New("sw")
	vars := url.Values{
		"first_name":       {"geofrey"},
		"last_name":        {"enrnest"},
		"middle_name":      {"gernest"},
		"email":            {"gernest@mawazo.com"},
		"password":         {"kilimahewa"},
		"confirm_password": {"kilimahewa"},
		"gender":           {"1"},
		"birth_date":       {"2 January, 1980"},
	}
	req1, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form1 := f.RegisterForm()(req1)
	if !form1.IsValid() {
		t.Errorf("validating form %v", form1.Errors())
	}

	// Case wrong email
	vars.Set("email", "kilimanjaro")
	req2, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form2 := f.RegisterForm()(req2)
	if form2.IsValid() {
		t.Error("expected validation error")
	}

	// Case wrong short password
	vars.Set("email", "kilimanjaro@example.com")
	vars.Set("password", "short")
	req3, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Errorf("creating new request %v", err)
	}
	req3.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form3 := f.RegisterForm()(req3)
	if form3.IsValid() {
		t.Error("expected validation error")
	}
	errs := form3.Errors()
	passErrs := errs.Get("password")
	if !strings.Contains(passErrs[0], f.tr.T(msgMinLength, 6)) {
		t.Errorf("expected %v to contain %v", errs, f.tr.T(msgMinLength))
	}

	// pass & confirm pass mismatch
	confirmErr := errs.Get("confirm_password")[0]
	expectMsg := f.tr.T(msgEqual, "confirm_password", "password")
	if expectMsg != confirmErr {
		t.Errorf("expected %s got %s", expectMsg, confirmErr)
	}
}
