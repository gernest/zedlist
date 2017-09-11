// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package flash implements flash messages, using gorilla/session.
package flash

import (
	"encoding/gob"

	"github.com/labstack/echo"
	"github.com/zedio/zedlist/modules/log"
	"github.com/zedio/zedlist/modules/session"
	"github.com/zedio/zedlist/modules/settings"
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

func (f *Flash) Class() string {
	switch f.Kind {
	case settings.FlashErr:
		return "error"
	case settings.FlashSuccess:
		return "success"
	case settings.FlashWarn:
		return "warning"
	default:
		return ""
	}
}

// Flashes is a collection of flash messages
type Flashes []*Flash

// GetFlashes retieves all flash messages found in a cookie session associated with ctx..
//
// BUG multipleflash messages are not propery set. the flashes contains only the first
// message to be set.
func GetFlashes(ctx echo.Context, key string) Flashes {
	ss, err := store.Get(ctx.Request(), settings.App.Session.Flash)
	if err != nil {
		//log.Error(nil, err)
	}
	if v, ok := ss.Values[key]; ok {
		delete(ss.Values, key)
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
	f := GetFlashes(ctx, settings.FlashKey)
	if f != nil {
		ctx.Set(settings.FlashKey, f)
	}
	fctx := GetFlashes(ctx, settings.FlashCtxKey)
	if fctx != nil {
		m := make(map[string]string)
		for _, v := range fctx {
			m[v.Kind] = v.Message
		}
		ctx.Set(settings.FlashCtxKey, m)
	}
}

//Flasher tracks flash messages
type Flasher struct {
	f   Flashes
	ctx Flashes
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

func (f *Flasher) AddCtx(kind, message string) {
	fl := &Flash{kind, message}
	f.ctx = append(f.ctx, fl)
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
	var flashes, flashesCtx Flashes
	if v, ok := ss.Values[settings.FlashKey]; ok {
		flashes = v.(Flashes)
	}

	if v, ok := ss.Values[settings.FlashCtxKey]; ok {
		flashesCtx = v.(Flashes)
	}
	ss.Values[settings.FlashKey] = append(flashes, f.f...)
	ss.Values[settings.FlashCtxKey] = append(flashesCtx, f.ctx...)
	err = ss.Save(ctx.Request(), ctx.Response())
	if err != nil {
		return err
	}
	return nil
}
