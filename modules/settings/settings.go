// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package settings contains all settings and configurations for zedlist.
package settings

import (
	"strings"

	"github.com/koding/multiconfig"
)

var (
	//App contains application specific variables
	App *Config

	// DocsPath is the relative directory for documents.

	// CodecsKeyPair is a slice of key pairs used to secure cookies
	CodecsKeyPair [][]byte
)

const (
	// DefaultConfig is the default configuration file
	DefaultConfig = "config.toml"

	// DataKey is the key used to store template specific context object, inside the
	// *echo.Context
	DataKey = "_zeddata"

	// LangSessionName is the name for lanuage session
	LangSessionName = "_lang"

	// LangSessionKey is the key for default language in language session
	LangSessionKey = "target"

	// LangSessionSecret is the secret for language session
	LangSessionSecret = "so-secret"

	// LangDataKey s the key for language context passed to templates
	LangDataKey = "ActiveLang"

	//JobsListKey the key in the data context containing a slice of jobs that is passed
	// to templates
	JobsListKey = "jobs"

	// PageTitleKey the key in the data context containing the pageTitle
	PageTitle = "PageTitle"

	// MaxListLimit the maximun number of items in a job list.
	MaxListLimit = 20

	//JobsFound the total number of jobs found
	JobsFound = "JobsFound"

	// RegionsListKey context key for regions list
	RegionsListKey = "regions"

	DocsPath = "static/docs"

	// DocsIndexPage the home page for documents
	DocsIndexPage = "home.md"

	// FlashSuccess is the context key for success flash messages
	FlashSuccess = "FlashSuccess"

	// FlashWarn is a context key for warning flash messages
	FlashWarn = "FlashWarn"

	// FlashErr is a context key for flash error message
	FlashErr = "FlashError"

	// FlashKey is the context key, passed to templates which holds lash messages
	FlashKey = "Flash"

	FlashCtxKey = "FlashCtx"

	SupportedLangs = "SupportedLangs"
	CurrentPath    = "CurrentPath"
)

func init() {
	App = NewConfig()
	configure()
}

// Config contains application specific configurations.
type Config struct {
	Name            string `default:"zedlist"`
	Port            int    `default:"8090"`
	AppURL          string `default:"http://localhost:8090"`
	MinimumAge      int    `default:"18"`
	BirthDateFormat string `default:"2 January, 2006"`
	DefaultLang     string `default:"en"`
	Session         Session

	// DBDialect is the type of database used.
	// supported values are mysql,postgres,foundation and sqlite.
	// default value is postgres.
	DBDialect string `default:"postgres"`

	// DBConn is the database connection string. This defaults to a connection to a postgres
	// database with the following credentials
	//
	// User	: postgres
	// Password	: postgres
	// Host	: localhost
	// sslmode	: disable
	DBConn string `default:"postgres://postgres:postgres@localhost/zedlist?sslmode=disable"`
}

//flash messages keys
const (
	FlashAccountCreate       = "flash_account_create"
	FlashAccountCreateFailed = "flash_account_create_fail"
	FlashLoginSuccess        = "flash_login_success"
	FlashLoginErr            = "flash_login_failed"
	FlashNotAuthorized       = "flash_unauthorized"
	FlashUnknownAccount      = "flash_unknown_account"
	FlashCreateJobSuccess    = "flash_create_new_job_success"
	FlashFailedNewJob        = "flash_failed_creating_new_job"
)

//keys for error messages
const (
	ErrorResourceNotFound = "error_resource_not_found"
	ErrorInternalServer   = "error_internal_server"
	ErrBadRequest         = "error_bad_request"
	Message               = "Message"
	ErrInvalidForm        = "error_ivalid_form"
)

// Session configurations for sessions
type Session struct {
	Name     string `default:"_zlst_"`
	Path     string `default:"/"`
	Domain   string
	MaxAge   int `default:"2592000"`
	Secure   bool
	HTTPOnly bool

	// Flash is the session name for flash messages
	Flash string `default:"_flash"`

	// Lang is the session name for languages
	Lang string `default:"_lang"`

	// KeyPair for secure cookie its a comma separates strings of keys.
	KeyPair string `default:"ePAPW9vJv7gHoftvQTyNj5VkWB52mlza,N8SmpJ00aSpepNrKoyYxmAJhwVuKEWZD"`
}

func configure() {
	// set codecs key pairs
	pairs := strings.Split(App.Session.KeyPair, ",")
	for _, key := range pairs {
		CodecsKeyPair = append(CodecsKeyPair, []byte(key))
	}
}

func newConfigLoader() *multiconfig.DefaultLoader {
	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{},
	)
	d := &multiconfig.DefaultLoader{}
	d.Loader = loader
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	return d
}

// NewConfig returns a new loaded Config.
func NewConfig() *Config {
	l := newConfigLoader()
	cfg := &Config{}
	l.MustLoad(cfg)
	return cfg
}
