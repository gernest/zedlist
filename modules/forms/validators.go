// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forms

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gernest/gforms"
)

var birthDateFormat = "2 January, 2006"

// EqualValidator checks if the two fields are equal. The to attribute is the name of
// the field whose value must be equal to the current field
type EqualValidator struct {
	To      string
	Message string
	gforms.Validator
}

// Validate  checks if the given field is egual to the field in the to attribute
func (vl EqualValidator) Validate(fi *gforms.FieldInstance, fo *gforms.FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.String || v.Value == "" {
		return nil
	}
	fi2, ok := fo.GetField(vl.To)
	if ok {
		v2 := fi2.GetV()
		if v.Value != v2.Value {
			return fmt.Errorf(vl.Message, fi.GetName(), fi2.GetName())
		}
	}
	return nil
}

// BirthDateValidator validates the birth date, handy to keep minors offsite
type BirthDateValidator struct {
	Limit   int
	Message string
	gforms.Validator
}

// Validate checks if the given field instance esceeds the Limit attribute
func (vl BirthDateValidator) Validate(fi *gforms.FieldInstance, fo *gforms.FormInstance) error {
	v := fi.V
	now := time.Now()
	var rerr = errors.New(vl.Message)

	if v.IsNil {
		return nil
	}
	iv := v.Value.(time.Time)
	if now.Year()-iv.Year() < vl.Limit {
		return rerr
	}
	return nil
}
