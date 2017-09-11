// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/zedio/zedlist/modules/settings"
)

func TestDB_Connect(t *testing.T) {
	d := NewDB(settings.NewConfig())
	db, err := d.Connect()
	if err != nil {
		t.Errorf("opening connection %v", err)
	}
	err = db.DB().Ping()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
}
