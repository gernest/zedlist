// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package flash implements flash messages, using gorilla/session.
package flash

import (
	"encoding/gob"

	"github.com/gernest/zedlist/modules/log"
	"github.com/gernest/zedlist/modules/session"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/labstack/echo"
)

func init() {
	gob.Register(&Flash{})
	gob.Register(Flashes{})
}

var store = session.New()

// Flash implements flash messages, like ones in gorilla/sessions
type Flash struct {
	Kind    string
	Message string
}

// Flashes is a collection of flash messages
type Flashes []*Flash

// GetFlashes retieves all flash messages found in a cookie session associated with ctx..
//
// BUG multipleflash messages are not propery set. the flashes contains only the first
// message to be set.
func GetFlashes(ctx echo.Context) Flashes {
	ss, err := store.Get(ctx.Request(), settings.App.Session.Flash)
	if err != nil {
		//log.Error(nil, err)
	}
	if v, ok := ss.Values[settings.FlashKey]; ok {
		delete(ss.Values, settings.FlashKey)
		serr := ss.Save(ctx.Request(), ctx.Response())
		if serr != nil {
			log.Error(ctx, err)
		}
		return v.(Flashes)
	}
	return nil
}

// AddFlashToCtx takes flash messages stored in a cookie which is associated with the
// request found in ctx, and puts them inside the ctx object. The flash messages can then
// be retrived by calling ctx.Get(settings.FlashKey).
//
// NOTE When there are no flash messages then nothing is set.
func AddFlashToCtx(ctx echo.Context) {
	f := GetFlashes(ctx)
	if f != nil {
		ctx.Set(settings.FlashKey, f)
	}
}

//Flasher tracks flash messages
type Flasher struct {
	f Flashes
}

//New creates new flasher. This alllows accumulation of lash messages. To save the flash messages
//the Save method should be called explicitly.
func New() *Flasher {
	return &Flasher{}
}

// Add adds the flash message
func (f *Flasher) Add(kind, message string) {
	fl := &Flash{kind, message}
	f.f = append(f.f, fl)
}

// Success adds success flash message
func (f *Flasher) Success(msg string) {
	f.Add(settings.FlashSuccess, msg)
}

// Err adds error flash message
func (f *Flasher) Err(msg string) {
	f.Add(settings.FlashErr, msg)
}

// Warn adds warning flash message
func (f *Flasher) Warn(msg string) {
	f.Add(settings.FlashWarn, msg)
}

// Save saves flash messages to context
func (f *Flasher) Save(ctx echo.Context) error {
	ss, err := store.Get(ctx.Request(), settings.App.Session.Flash)
	if err != nil {
		//log.Error(nil, err)
	}
	var flashes Flashes
	if v, ok := ss.Values[settings.FlashKey]; ok {
		flashes = v.(Flashes)
	}
	ss.Values[settings.FlashKey] = append(flashes, f.f...)
	err = ss.Save(ctx.Request(), ctx.Response())
	if err != nil {
		return err
	}
	return nil
}
