// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package log provides logrus based logging for zedlist.
package log

import (
	"time"

	"github.com/labstack/echo"

	"github.com/Sirupsen/logrus"
)

// L is the global logger
var L *logrus.Logger

//var WL *WorkerLog

var maxTimeout = time.Second

func init() {
	//WL = NewWOrkerLog(work.W)
	L = New()
}

//type WorkerLog struct {
//	send   chan echo.Context
//	stop   chan bool
//	server *gowork.WorkServer
//}

//func NewWOrkerLog(ws *gowork.WorkServer) *WorkerLog {
//	w := &WorkerLog{
//		send:   make(chan echo.Context, 100),
//		server: ws,
//	}
//	w.init()
//	return w
//}

//func (wl *WorkerLog) init() {
//	wl.server.NewHandler("add_work", wl.onWorkAdd)
//}

//func (wl *WorkerLog) Run() {
//	defer func() {
//		msg := recover()
//		fmt.Printf("recovering from %v", msg)
//	}()

//	// add jobs to the queue.
//	go func() {
//	OUT:
//		for {
//			select {
//			case nwork := <-wl.send:
//				params := make(map[string]interface{})
//				params["ctx"] = nwork
//				newWork, err := gowork.CreateWork(params, int64(maxTimeout))
//				if err != nil {
//					logrus.Error(err)
//					return
//				}
//				wl.server.Add(newWork)
//			case <-wl.stop:
//				break OUT
//			}

//		}
//	}()
//}

//func (wl *WorkerLog) onWorkAdd(event *gowork.Event, params map[string]interface{}) {
//	go func() {
//	OUT:
//		for {
//			select {
//			case <-wl.stop:
//				break OUT
//			default:
//				if eCtx, ok := params["ctx"]; ok {
//					ctx := eCtx.(echo.Context)
//					ctxPtr := &ctx
//					loggerStruct := models.Logger{}
//					usr := ctxPtr.Get("User")
//					if usr != nil {
//						person := usr.(*models.Person)
//						loggerStruct.User = person.Email
//					}
//					if loggerStruct.User == "" {
//						loggerStruct.User = "John Doe"
//					}
//					level := ctx.Get("level").(logrus.Level)
//					switch level {
//					case logrus.PanicLevel:
//						loggerStruct.Level = PanicLevel
//					case logrus.FatalLevel:
//						loggerStruct.Level = FatalLevel
//					case logrus.ErrorLevel:
//						loggerStruct.Level = ErrorLevel
//					case logrus.WarnLevel:
//						loggerStruct.Level = WarnLevel
//					case logrus.InfoLevel:
//						loggerStruct.Level = InfoLevel
//					case logrus.DebugLevel:
//						loggerStruct.Level = DebugLevel
//					}
//					err := query.Create(loggerStruct)
//					if err != nil {
//						logrus.Error(err)
//					}
//					break OUT
//				}
//			}
//		}
//	}()
//}

//func (wl *WorkerLog) Push(ctx echo.Context) {
//	wl.send <- ctx
//}

// Error logs error messages
func Error(ctx echo.Context, args ...interface{}) {
	if ctx != nil {
		d := ctx
		L.WithField("ctx", d).Error(args...)
	} else {
		L.Error(args...)
	}

}

// Info logs info messages
func Info(ctx echo.Context, args ...interface{}) {
	if ctx != nil {
		d := ctx
		L.WithField("ctx", d).Info(args...)
	} else {
		L.Info(args...)
	}

}

// New creates a new logger
func New() *logrus.Logger {
	nlog := logrus.New()
	nlog.Formatter = &ZedFormatter{}
	return nlog
}
