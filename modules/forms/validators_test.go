// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gernest/gforms"
)

func TestBirthDateValidator(t *testing.T) {
	var now = time.Now()
	var yearsAgo = func(yrs int) time.Time {
		n := time.Now()
		nowAFter := n.AddDate(18, 1, 1)
		dur := nowAFter.Sub(n)
		return n.Add(-dur)

	}
	Form := gforms.DefineForm(gforms.NewFields(
		gforms.NewDateTimeField(
			"date",
			time.RFC822,
			gforms.Validators{
				BirthDateValidator{Limit: 18, Message: msgAge},
			},
		),
	))

	vars := url.Values{"date": {now.Format(time.RFC822)}}
	req1, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Error(err)
	}
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form := Form(req1)
	if form.IsValid() {
		t.Error("Expected some errors")
	}

	vars = url.Values{"date": {yearsAgo(18).Format(time.RFC822)}}
	req2, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Error(err)
	}
	req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form = Form(req2)
	if !form.IsValid() {
		t.Error(form.Errors())
	}
}

func TestEqualValidator(t *testing.T) {
	var msg = "%s should equal %s"
	eForm := gforms.DefineForm(gforms.NewFields(
		gforms.NewTextField(
			"first",
			gforms.Validators{
				gforms.Required(),
			},
		),
		gforms.NewTextField(
			"second",
			gforms.Validators{
				EqualValidator{To: "first", Message: msg},
			},
		),
	))

	vars := url.Values{
		"first":  {"hello"},
		"second": {"world"},
	}
	req1, err := http.NewRequest("POST", "/", strings.NewReader(vars.Encode()))
	if err != nil {
		t.Error(err)
	}
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	form := eForm(req1)
	if form.IsValid() {
		t.Error("expected errors")
	}
	errMsg := form.Errors().Get("second")[0]
	expMsg := fmt.Sprintf(msg, "second", "first")
	if errMsg != expMsg {
		t.Errorf(" expected %s got %s", expMsg, errMsg)
	}
}
