// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sanitize contains sanitization policies for zedlist.This adds protection against XSS
// attacks.
package sanitize

import (
	"github.com/microcosm-cc/bluemonday"
)

var strict *bluemonday.Policy

func init() {
	strict = bluemonday.StrictPolicy()
}

// Name sanitizes names
func Name(str string) string {
	return strict.Sanitize(str)
}

// Title sanitizes titles
func Title(str string) string {
	return strict.Sanitize(str)
}
