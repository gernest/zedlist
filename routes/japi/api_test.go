// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package japi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/Unknwon/com"
	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
)

var ts = testServer()
var client = &http.Client{Jar: getJar()}

func testServer() *httptest.Server {
	e := echo.New()
	a := e.Group("/api")
	a.Post("/jobs", CreateJob)
	a.Get("/jobs/:id", GetJob)
	a.Get("/jobs", GetIndex)
	a.Put("/jobs", UpdateJob)
	return httptest.NewServer(e)
}

func getJar() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Log
	}
	return jar
}

func cleanTables() {
	db.Conn.DropTableIfExists(&models.Job{})
	db.Conn.AutoMigrate(&models.Job{})
}

func closeTest() {
	db.Conn.Close()
}

//
//
//		JOB API
//
//
func TestCreateJob(t *testing.T) {
	job := &models.Job{
		Title:       "my job",
		Description: "whacko job",
	}
	pathURL := fmt.Sprintf("%s/api/jobs", ts.URL)
	rst := models.Job{}
	err := com.HttpPostJSON(client, pathURL, job, &rst)
	if err != nil {
		t.Errorf("creating a job %v", err)
	}
	if rst.Title != job.Title {
		t.Errorf("expected %s got %s", job.Title, rst.Title)
	}

}

func TestGetJob(t *testing.T) {
	pathURL := fmt.Sprintf("%s/api/jobs", ts.URL)
	jsonErr := models.JSONError{}
	err := com.HttpGetJSON(client, fmt.Sprintf("%s/1000", pathURL), &jsonErr)
	if err == nil {
		t.Error("expected resource not found")
	}
	job := &models.Job{
		Title:       "my job",
		Description: "whacko job",
	}
	q := db.Conn.Create(job)
	if q.Error != nil {
		t.Error(q.Error)
	}

	rst := &models.Job{}
	err = com.HttpGetJSON(client, fmt.Sprintf("%s/%d", pathURL, job.ID), rst)
	if err != nil {
		t.Errorf("getting record %v", err)
	}
	if rst.Title != job.Title {
		t.Errorf("expected %s got %s", job.Title, rst.Title)
	}
	if rst.ID != job.ID {
		t.Errorf("expected %d got %d", job.ID, rst.ID)
	}
}

func TestUpdateJob(t *testing.T) {
	job := &models.Job{
		Title:       "my job",
		Description: "whacko job",
	}
	pathURL := fmt.Sprintf("%s/api/jobs", ts.URL)
	rst := models.Job{}
	err := com.HttpPostJSON(client, pathURL, job, &rst)
	if err != nil {
		t.Errorf("creating a job %v", err)
	}
	if rst.Title != job.Title {
		t.Errorf("expected %s got %s", job.Title, rst.Title)
	}
	rst.Description = "a changed world"
	upRst := models.Job{}
	err = httpPutJSON(client, pathURL, &rst, upRst)
	if err != nil {
		t.Error(err)
	}
	check := models.Job{}
	q := db.Conn.Where(&models.Job{ID: rst.ID}).First(&check)
	if q.Error != nil {
		t.Error(q.Error)
	}

	if check.Description != rst.Description {
		t.Errorf("expected %s got %s", rst.Description, check.Description)
	}

}

func httpPut(client *http.Client, url string, header http.Header, body []byte) (io.ReadCloser, error) {
	return com.HttpCall(client, "PUT", url, header, bytes.NewBuffer(body))
}

func httpPutJSON(client *http.Client, url string, body, v interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	rc, err := httpPut(client, url, http.Header{"content-type": []string{"application/json"}}, data)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		return fmt.Errorf("JSON syntax error at %s", url)
	}
	return nil
}

func TestGetIndex(t *testing.T) {
	data := []struct {
		title string
	}{
		{"geernest"},
		{"mwanza"},
		{"Tanzania"},
	}

	// create all the records.
	for _, v := range data {
		j := models.Job{Title: v.title}
		q := db.Conn.Create(&j)

		if q.Error != nil {
			t.Errorf("creating a job record %v", q.Error)
		}
	}

	pathURL := fmt.Sprintf("%s/api/jobs", ts.URL)
	rst := []*models.Job{}
	err := com.HttpGetJSON(client, pathURL, &rst)
	if err != nil {
		t.Errorf("retrieving jobs index %v", err)
	}
	if len(rst) < len(data) {
		t.Errorf("expected %d got %d", len(data), len(rst))
	}

	// verify the order is in descending
	if rst[0].Title != data[2].title {
		t.Errorf("expected %s got %s", data[2].title, rst[0].Title)
	}
}

// this should be always at the bottom of this file.
func TestCLose(t *testing.T) {
	closeTest()
}
