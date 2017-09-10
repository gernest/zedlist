// Copyright 2015 Geofrey Ernest <geofreyernest@live.com>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package i18n is a translation library.
package i18n

import (
	"github.com/gernest/gt"
	"github.com/gernest/zedlist/modules/settings"
)

// Tr contains translations.
var Tr = &gt.Build{
	Origin: "en",
	Target: "sw",
	Index: gt.Strings{
		"home-btn": {
			"en": "Home",
			"sw": "Nyumbani",
		},
		"jobs-btn": {
			"en": "Jobs",
			"sw": "Ajira",
		},
		"job": {
			"en": "Job",
			"sw": "Ajira",
		},
		"resume": {
			"en": "Resume",
			"sw": "Rejea",
		},
		"help-btn": {
			"en": "Help",
			"sw": "Msaada",
		},
		"contact-btn": {
			"en": "contact us",
			"sw": "wasiliana nasi",
		},
		"deadline": {
			"en": "deadline",
			"sw": "mwisho",
		},
		"regions": {
			"en": "regions",
			"sw": "mikoa",
		},
		"apply-btn": {
			"en": " apply now",
			"sw": "omba sasa",
		},
		"login": {
			"en": "Login",
			"sw": "Ingia",
		},
		"logout": {
			"en": "Logout",
			"sw": "Jitoe",
		},
		"register": {
			"en": "Register",
			"sw": "Iiunge",
		},
		"message_required": {
			"en": "this field should not be empty",
			"sw": "hili eneo halitakiwi kuachwa wazi",
		},
		"message_min_length": {
			"en": "this field should be at least %d characters",
			"sw": "namba ya siri inatakiwa kuanzia herufi %d na kuendelea",
		},
		"message_email": {
			"en": "incorrect email, should of the form example@examples.com",
			"sw": "barua pepe sio sahihi. mfano example@example.com",
		},
		"message_age": {
			"en": "age should be more than %d years",
			"sw": "umri unatakiwa uwe zaidi ya miaka %d",
		},
		"message_equal": {
			"en": "%s should be equal to %s",
			"sw": "%s inatakiwa iwe sawa na %s",
		},
		"documents": {
			"en": "documents",
			"sw": "makala",
		},

		"issued_by": {
			"en": "issued by",
			"sw": "imetolewa na",
		},
		"valid_name_msg": {
			"en": "field %s should conatin letters a-zA-Z e.g baba",
			"sw": "eneo %s linatakiwa liwe na herufi to a-zA-Z mfano baba",
		},
		"given_name": {
			"en": "first name",
			"sw": "jina la kwanza",
		},
		"family_name": {
			"en": "family name",
			"sw": "jina la ukoo",
		},
		"middle_name": {
			"en": "middle name",
			"sw": "jina la baba",
		},
		"username": {
			"en": "Username",
			"sw": "Jina la akaunti",
		},
		"username_or_email": {
			"en": "Username or Email",
			"sw": "Jina la akaunti au barua pepe",
		},
		"password": {
			"en": "Password",
			"sw": "Neno la siri",
		},
		"confirm_password": {
			"en": "Confirm Password",
			"sw": "Thibitisha neno la siri",
		},
		"email": {
			"en": "Email",
			"sw": "Barua pepe",
		},
		"delete_account": {
			"en": "Delete account",
			"sw": "Funga akaunti",
		},
		"register_form_title": {
			"en": "Create a new account",
			"sw": "Tengeneza akaunti mpya",
		},
		"login_form_title": {
			"en": "Login to your account",
			"sw": "Ingia kwenye akaunti yako",
		},
		"new": {
			"en": "New",
			"sw": "Mpya",
		},
		"catch_phrase": {
			"en": "Get the right person for the Job",
			"sw": "Pata mtu sahihi kwa kazi sahihi",
		},
		"catch_phrase_sub": {
			"en": "A humble self hosted recruitment service",
			"sw": "Mfumo bora wa ajira",
		},
		"new_job_title": {
			"en": "Create a Job posting",
			"sw": "Tengeneza ajira mpya",
		},
		"job_title": {
			"en": "Title",
			"sw": "Kichwa cha habari",
		},
		"job_desc": {
			"en": "Description",
			"sw": "Maelezo",
		},
		"create": {
			"en": "Create",
			"sw": "Tengeneza",
		},
		"update": {
			"en": "Update",
			"sw": "Sahihisha",
		},
		"update_job_title": {
			"en": "Update a job posting",
			"sw": "Sahihisha habari za ajira",
		},
		"settings": {
			"en": "Settings",
			"sw": "Settings",
		},
	},
}

func init() {
	Tr.AddIndex(flash())
	Tr.AddIndex(errorMessages())
	Tr.Init()
}

// CloneLang returns a copy of translations.
func CloneLang() *gt.Build {
	return Tr
}

// flash message translation indices
func flash() gt.Strings {
	return gt.Strings{
		settings.FlashAccountCreate: {
			"en": "congard, your account has been successful created",
			"sw": "hongera, akaunti yako imefanikiwa kutengenezwa",
		},
		settings.FlashAccountCreateFailed: {
			"en": "sorry, we can't create you account please try again later",
			"sw": "samahaani, tumeshindwa kutengeneza akaunti yako, jaribu tena baadae",
		},
		settings.FlashLoginSuccess: {
			"en": "welcome back",
			"sw": "karibu ",
		},
		settings.FlashLoginErr: {
			"en": "there was a problem encountered, please check the details and try again",
			"sw": "kuna majanga mkuu, jaribu kupitia maelezo ya fomu na ujaribu tena",
		},
		settings.FlashNotAuthorized: {
			"en": "not authorized",
			"sw": "ihauna ruhusa",
		},
		settings.FlashUnknownAccount: {
			"en": "unknown account",
			"sw": "ihauna ruhusa",
		},
		settings.FlashCreateJobSuccess: {
			"en": "A new job was successful created",
			"sw": "Ajia mpya imetangenezwa",
		},
		settings.FlashFailedNewJob: {
			"en": "sorry! We failed to create the new job",
			"sw": "samahani tumeshindwa kutengeneza ajira mpya",
		},
		settings.FlashSuccessUpdate: {
			"en": "succesful updated job listing",
			"sw": "umefanikiwa kusahihisha habari za ajira",
		},
		settings.FlashSuccessDelete: {
			"en": "succesful deleted a job listing",
			"sw": "umefanikiwa kufuta ajira",
		},
	}
}

func errorMessages() gt.Strings {
	return gt.Strings{
		settings.ErrorResourceNotFound: {
			"en": "Ooops! we cant find the resource you asked for",
			"sw": "Oops! Tumeshindwa kupata ulichoulizia",
		},
		settings.ErrorInternalServer: {
			"en": "Opps! we had a problem",
			"sw": "Oops! kuna hitilafu",
		},
		settings.ErrBadRequest: {
			"en": "Opps! bad request",
			"sw": "Oops! bad request",
		},
		settings.ErrInvalidForm: {
			"en": "something is wron with the submitted form",
			"sw": "Kuna matatizo kwenye fomu",
		},
	}
}
