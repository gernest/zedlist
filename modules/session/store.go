// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package session is  postgresql based gorilla/sessions backend.
*/
package session

import (
	"encoding/base32"
	"net/http"
	"strings"
	"time"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/query"

	"github.com/gernest/zedlist/modules/settings"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// PGStore is the session storage implementation for gorilla/sessions using postgresql
type PGStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
}

// New initializes a new PGStore, using key pairs found in settings configuration
func New() *PGStore {
	return NewPGStore(settings.CodecsKeyPair...)
}

// NewPGStore initillizes PGStore with the given keyPairs
func NewPGStore(keyPairs ...[]byte) *PGStore {
	dbStore := &PGStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   settings.App.Session.Path,
			MaxAge: settings.App.Session.MaxAge,
		},
	}
	return dbStore
}

// Get fetches a session for a given name after it has been added to the registry.
func (db *PGStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(db, name)
}

// New returns a new session
func (db *PGStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(db, name)
	opts := *db.Options
	session.Options = &(opts)
	session.IsNew = true

	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, db.Codecs...)
		if err == nil {
			err = db.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save saves the session into a postgresql database
func (db *PGStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Set delete if max-age is < 0
	if session.Options.MaxAge < 0 {
		if err := db.Delete(r, w, session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Generate a random session ID key suitable for storage in the DB
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}

	if err := db.save(session); err != nil {
		return err
	}

	// Keep the session ID key in a cookie so it can be looked up in DB later.
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, db.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

//load fetches a session by ID from the database and decodes its content into session.Values
func (db *PGStore) load(session *sessions.Session) error {
	s, err := query.GetSessionByKey(session.ID)
	if err != nil {
		return err
	}
	return securecookie.DecodeMulti(session.Name(), string(s.Data),
		&session.Values, db.Codecs...)
}

func (db *PGStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		db.Codecs...)

	if err != nil {
		return err
	}

	var expiresOn time.Time

	exOn := session.Values["expires_on"]

	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
		if expiresOn.Sub(time.Now().Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
			expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
		}
	}
	s := &models.Session{
		Key:       session.ID,
		Data:      encoded,
		ExpiresOn: expiresOn,
	}
	if session.IsNew {
		return query.Create(s)
	}
	return query.UpdateSession(s)
}

func (db *PGStore) destroy(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	options := *db.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	for k := range session.Values {
		delete(session.Values, k)
	}
	return query.DeleteSession(session.ID)
}

// Delete deletes session.
func (db *PGStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return db.destroy(r, w, session)
}
