package models

import "git/ymoldabe/forum/validator"

const (
	ClientID     = "178180018520-gtopcrp8feslt8gbm19e7lmhrls7g9ai.apps.googleusercontent.com"
	ClientSecret = "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"
)

const (
	GoogleAuthURL = "https://accounts.google.com/o/oauth2/auth"

	GoogleRedirectUrl = "http://localhost:8081/auth/callback"

// googleTokenURL    = "https://accounts.google.com/o/oauth2/token"
// googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type UserSignUp struct {
	Name         string
	Email        string
	Password1    string
	Password2    string
	HashPassword string
	validator.Validator
}

type UserSignIn struct {
	Id           int
	Email        string
	Password     string
	HashPassword string
	validator.Validator
}
