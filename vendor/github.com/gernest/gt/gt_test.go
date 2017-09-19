// Copyright (c) 2013 Melvin Tercan, https://github.com/melvinmt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
// persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package gt

import (
	"testing"
)

func TestBuild(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"en":    "Welcome to %s#title, %s#name!",
				"es-LA": "¡Bienvenido a %s#title, %s#title!",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	_ = g
}

func TestNoTargetString(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"en": "Welcome to %s#title, %s#name!",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not Welcome to gt, Melvin!", tr)
			t.Error(err)
		}
		if err == nil {
			t.Error("Should return: could not find error")
		}
	}()
	func() {
		t.Log("Test string")
		s := "Welcome to %s#title, %s#name!"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not Welcome to gt, Melvin!", tr)
		}
		if err == nil {
			t.Error("Should return: could not find error")
		}
	}()
}

func TestNoOriginString(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"es-LA": "¡Bienvenido a %s#title, %s#title!",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s)
		if s != tr {
			t.Errorf("Returned string '%s' is not the same as input string '%s'", tr, s)
		}
		if err == nil {
			t.Error("Should return: could not find error")
		}
	}()
	func() {
		t.Log("Test string")
		s := "Welcome to %s#title, %s#name!"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not the same as Welcome to gt, Melvin!", tr)
		}
		if err == nil {
			t.Error("Should return: could not find error")
		}
	}()
}

func TestNoIndex(t *testing.T) {
	g := &Build{
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s)
		if s != tr {
			t.Errorf("Returned string '%s' is not the same as input string '%s'", tr, s)
		}
		if err == nil {
			t.Error("Should return: target or origin is not set")
		}
	}()
	func() {
		t.Log("Test string")
		s := "Welcome to %s#title, %s#name!"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not the same as Welcome to gt, Melvin!", tr)
		}
		if err == nil {
			t.Error("Should return: target or origin is not set")
		}
	}()
}

func TestNoOrigin(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"en":    "Welcome to %s#title, %s#name!",
				"es-LA": "¡Bienvenido a %s#title, %s#title!",
			},
		},
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s)
		if s != tr {
			t.Errorf("Returned string '%s' is not the same as input string '%s'", tr, s)
		}
		if err == nil {
			t.Error(err)
			t.Error("Should return: target or origin is not set")
		}
	}()
	func() {
		t.Log("Test string")
		s := "Welcome to %s#title, %s#name!"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not the same as Welcome to gt, Melvin!", tr)
		}
		if err == nil {
			t.Error("Should return: target or origin is not set")
		}
	}()
}

func TestNoTarget(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"en":    "Welcome to %s#title, %s#name!",
				"es-LA": "¡Bienvenido a %s#title, %s#title!",
			},
		},
		Origin: "en",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s)
		if tr != "Welcome to %s(MISSING), %s(MISSING)!" {
			t.Errorf("Returned string '%s' is not the same as Welcome to %s(MISSING), %s(MISSING)!", tr)
		}
		if err == nil {
			t.Error("Should return: target or origin is not set")
		}
	}()
	func() {
		t.Log("Test string")
		s := "Welcome to %s#title, %s#name!"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "Welcome to gt, Melvin!" {
			t.Errorf("Returned string '%s' is not the same as Welcome to gt, Melvin!", tr)
		}
		if err == nil {
			t.Error("Should return: target or origin is not set")
		}
	}()
}

func TestNoArgsTranslation(t *testing.T) {
	g := &Build{
		Index: Strings{
			"homepage-greeting": {
				"en":    "Welcome to gt!",
				"es-LA": "¡Bienvenido a gt!",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test key")
		s := "homepage-greeting"
		tr, err := g.Translate(s)
		if tr != "¡Bienvenido a gt!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test key")
		s := "Welcome to gt!"
		tr, err := g.Translate(s)
		if tr != "¡Bienvenido a gt!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestArgsTranslationValidSyntax(t *testing.T) {
	g := &Build{
		Index: Strings{
			"single": {
				"en":    "Welcome to %s!",
				"es-LA": "¡Bienvenido a %s!",
			},
			"multiple": {
				"en":    "Welcome to %s, %s!",
				"es-LA": "¡Bienvenido a %s, %s!",
			},
			"tags": {
				"en":    "Welcome to %s#title-web_site, %s#name2!",
				"es-LA": "¡Bienvenido a %s#title-web_site, %s#name2!",
			},
			"swap": {
				"en":    "Welcome to %s#title-web_site, %s#name2!",
				"es-LA": "¡Bienvenido a %s#name2, %s#title-web_site!",
			},
			"duplicate-same-order": {
				"en":    "Welcome to %s#title-web_site, %s#name2! Your name is the same: %s#name2",
				"es-LA": "¡Bienvenido a %s#title-web_site, %s#name2! Tu nombre es mismo: %s#name2",
			},
			"all-verbs": { // %b is duplicate but because it's in the same order (no swapping), no error is thrown
				"en":    "Welcome to %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %b, %e, %f, %g, %G, %s, %+q, % d, %12d, %1.2f!",
				"es-LA": "¡Bienvenido a %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %b, %e, %f, %g, %G, %s, %+q, % d, %12d, %1.2f!",
			},
			"all-verbs-tags": {
				"en":    "Welcome to %v#t1, %#v#t2, %T#t3, %%, %t#t4, %b#t5, %c#t6, %d#t7, %o#t8, %q#t9, %x#t10, %X#t11, %U#t12, %b#t13, %e#t14, %f#t15, %g#t16, %G#t17, %s#t18, %+q#t19, % d#t20, %12d#t21, %1.2f#t22!",
				"es-LA": "¡Bienvenido a %v#t1, %#v#t2, %T#t3, %%, %t#t4, %b#t5, %c#t6, %d#t7, %o#t8, %q#t9, %x#t10, %X#t11, %U#t12, %b#t13, %e#t14, %f#t15, %g#t16, %G#t17, %s#t18, %+q#t19, % d#t20, %12d#t21, %1.2f#t22!",
			},
			"all-verbs-different-order": { // notice I've removed a duplicate %b because that would cause a non unique swapping error
				"en":    "Welcome to %v, %#v, %T, %%, %t, %c, %d, %o, %q, %x, %X, %U, %b, %e, %f, %g, %G, %s, %+q, % d, %12d, %1.2f!",
				"es-LA": "¡Bienvenido a %b, %e, %f, %g, %G, %s, %+q, % d, %12d, %1.2f, %v, %#v, %T, %%, %t, %c, %d, %o, %q, %x, %X, %U!",
			},
			"all-verbs-tags-different-order": { // also notice we don't have this error with tags, because every tag provides an unique combination
				"en":    "Welcome to %v#t1, %#v#t2, %T#t3, %%, %t#t4, %b#t5, %c#t6, %d#t7, %o#t8, %q#t9, %x#t10, %X#t11, %U#t12, %b#t13, %e#t14, %f#t15, %g#t16, %G#t17, %s#t18, %+q#t19, % d#t20, %12d#t21, %1.2f#t22!",
				"es-LA": "¡Bienvenido a %e#t14, %f#t15, %g#t16, %G#t17, %s#t18, %+q#t19, % d#t20, %12d#t21, %1.2f#t22, %v#t1, %#v#t2, %T#t3, %%, %t#t4, %b#t5, %c#t6, %d#t7, %o#t8, %q#t9, %x#t10, %X#t11, %U#t12, %b#t13!",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test single %s")
		s := "single"
		tr, err := g.Translate(s, "gt")
		if tr != "¡Bienvenido a gt!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test multiple %s")
		s := "multiple"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "¡Bienvenido a gt, Melvin!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test tags")
		s := "tags"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "¡Bienvenido a gt, Melvin!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()

	func() {
		t.Log("Test swap tags")
		s := "swap"
		tr, err := g.Translate(s, "gt", "Melvin")
		if tr != "¡Bienvenido a Melvin, gt!" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test duplicate verbs in same order")
		s := "duplicate-same-order"
		tr, err := g.Translate(s, "gt", "Melvin", "Melvin")
		if tr != "¡Bienvenido a gt, Melvin! Tu nombre es mismo: Melvin" {
			t.Errorf("Returned string '%s' is not a translation of '%s'", tr, s)
		}
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test all types of verbs")
		s := "all-verbs"
		// "¡Bienvenido a %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %e, %f, %g, %G, %s, %+q, % d, %5d, %1.2f!",
		_, err := g.Translate(s, "s1", "s2", "s3", true, 1, 2, 3, 4, 5, 6, 7, 8, 1.1, 1.2, 1.3, 1.4, 1.5, "s4", 8, 9, 10, 1.6)
		// its almost impossible to test for str match so just checking for errors:
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test all types of verbs w/ tags")
		s := "all-verbs-tags"
		// "¡Bienvenido a %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %e, %f, %g, %G, %s, %+q, % d, %5d, %1.2f!",
		_, err := g.Translate(s, "s1", "s2", "s3", true, 1, 2, 3, 4, 5, 6, 7, 8, 1.1, 1.2, 1.3, 1.4, 1.5, "s4", 8, 9, 10, 1.6)
		// its almost impossible to test for str match so just checking for errors:
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test all types of verbs in different order")
		s := "all-verbs-different-order"
		// "¡Bienvenido a %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %e, %f, %g, %G, %s, %+q, % d, %5d, %1.2f!",
		_, err := g.Translate(s, 1.1, 1.2, 1.3, 1.4, 1.5, "s4", 8, 9, 10, 1.6, "s1", "s2", "s3", true, 1, 3, 4, 5, 6, 7, 8)
		// its almost impossible to test for str match so just checking for errors:
		if err != nil {
			t.Error(err)
		}
	}()
	func() {
		t.Log("Test all types of verbs w/ tags in different order")
		s := "all-verbs-tags-different-order"
		// "¡Bienvenido a %v, %#v, %T, %%, %t, %b, %c, %d, %o, %q, %x, %X, %U, %e, %f, %g, %G, %s, %+q, % d, %5d, %1.2f!",
		_, err := g.Translate(s, 1.1, 1.2, 1.3, 1.4, 1.5, "s4", 8, 9, 10, 1.6, "s1", "s2", "s3", true, 1, 2, 3, 4, 5, 6, 7, 8)
		// its almost impossible to test for str match so just checking for errors:
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestArgsTranslationInvalidSyntax(t *testing.T) {
	g := &Build{
		Index: Strings{
			"duplicate-different-order": {
				"en":    "Welcome to %s#title-web_site, %s#name2! Your name is the same: %s#name2",
				"es-LA": "¡Bienvenido a %s#name2, %s#title-web_site! Tu nombre es mismo: %s#name2",
			},
		},
		Origin: "en",
		Target: "es-LA",
	}
	g.Init()
	func() {
		t.Log("Test duplicate verbs in different order")
		s := "Welcome to %s#title-web_site, %s#name2! Your name is the same: %s#name2"
		tr, err := g.Translate(s, "gt", "Melvin", "Melvin")
		if tr != "Welcome to gt, Melvin! Your name is the same: Melvin" {
			t.Errorf("Returned string '%s' is not the same as Welcome to gt, Melvin! Your name is the same: Melvin", tr)
		}
		if err == nil {
			t.Error("Should return Verbs have to be swapped but are not unique error")
		}
	}()
}

func TestSetTargetOrOrigin(t *testing.T) {
	g := &Build{
		Index: Strings{
			"duplicate-different-order": {
				"en":    "Welcome to %s#title-web_site, %s#name2! Your name is the same: %s#name2",
				"es-LA": "¡Bienvenido a %s#name2, %s#title-web_site! Tu nombre es mismo: %s#name2",
			},
		},
	}
	g.Init()
	func() {
		t.Log("setting Target")
		g.Target = "en"
		g.SetTarget("es")
	}()
	func() {
		t.Log("setting Origin")
		g.Origin = "en"
		g.SetOrigin("es")
	}()
}
