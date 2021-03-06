// Copyright 2015-2017 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package japi is jobs API handlers for zedlist.
package japi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/zedio/zedlist/models"
	"github.com/zedio/zedlist/modules/db"
	"github.com/zedio/zedlist/modules/query"
	"github.com/zedio/zedlist/modules/utils"

	"github.com/labstack/echo"
)

var (
	errNotFound = errors.New("not found")
)

// CreateJob creates a new job record
func CreateJob(ctx echo.Context) error {
	job := &models.Job{}
	err := unmarshalToJSON(ctx, job)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}

	// sanitize before saving.
	job.Sanitize()

	err = query.CreateJob(db.Conn, job)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}
	return ctx.JSON(http.StatusOK, job)
}

// unmarhsalls the request body to val(in json format)).
func unmarshalToJSON(ctx echo.Context, val interface{}) error {
	req := ctx.Request()
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf.Bytes(), val)
	if err != nil {
		return err
	}
	return nil
}

// GetJob retrieves a job by ID
func GetJob(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}
	job, err := query.GetJobByID(db.Conn, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}

	// Sanitize before rendering
	job.Sanitize()

	return ctx.JSON(http.StatusOK, job)
}

// GetIndex retrieves all jobs.
func GetIndex(ctx echo.Context) error {
	jobs, err := query.GetALLJobs(db.Conn)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}
	return ctx.JSON(http.StatusOK, jobs)
}

// UpdateJob updates a job record
func UpdateJob(ctx echo.Context) error {
	job := &models.Job{}
	err := unmarshalToJSON(ctx, job)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}

	// sanitize before saving
	job.Sanitize()

	err = query.Update(db.Conn, job)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(err.Error()))
	}
	return ctx.JSON(http.StatusOK, job)
}
