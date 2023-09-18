package models

import (
	"git/ymoldabe/forum/validator"
	"time"

	"github.com/google/uuid"
)

const (
	GoogleClientID     = "178180018520-4spc30vqhq73nkan357nv7ldff4c6rf0.apps.googleusercontent.com"
	GoogleClientSecret = "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"
	GithubClientID     = "c40bb713a6ba597b0176"
	GithubClientSecret = "0e5fc1a85b8a3ebe6b2b04cc823d2a5050fe3e21"
)

const (
	GoogleAuthURL = "https://accounts.google.com/o/oauth2/auth"

	GoogleRedirectUrl = "http://localhost:8081/auth/google/callback"

	GitHubAuthURL = "https://github.com/login/oauth/authorize"

	GithubRedirectUrl = "http://localhost:8081/auth/github/callback"
)

const (
	ProviderGoogle = "google_user"
	ProviderDefult = "default"
	ProviderGit    = "github_user"
)

type UserSignUp struct {
	Name         string
	Email        string
	Password1    string
	Password2    string
	HashPassword string
	Provider     string
	validator.Validator
}

type UserSignIn struct {
	Id           int
	Email        string
	Password     string
	HashPassword string
	Provider     string
	validator.Validator
}

type GoogleLoginUserData struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      string
	Photo     string
	Verified  bool
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
	validator.Validator
}

type GithubLoginUserData struct {
	Id        int
	UserName  string `json:"login"`
	Password  string
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Role      string
	Provider  string
	validator.Validator
}
