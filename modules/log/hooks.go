// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

//import (
//	"fmt"
//	"github.com/labstack/echo"

//	"github.com/Sirupsen/logrus"
//)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

//const (
//	LoggerSecret = "32 bit secret"
//)

////PGHook is the postgres hook
//type PGHook struct{}

//func (p *PGHook) Fire(entry *logrus.Entry) error {
//	if eCtx, ok := entry.Data["ctx"]; ok {
//		ctx := eCtx.(echo.Context)
//		ctx.Set("level", entry.Level)
//		WL.Push(ctx)
//		return nil
//	}
//	return fmt.Errorf("bad entry %s", entry.Message)
//}

//func (p *PGHook) Levels() []logrus.Level {
//	return []logrus.Level{
//		logrus.ErrorLevel,
//		logrus.InfoLevel,
//	}
//}
