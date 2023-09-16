package store_auth

import (
	"time"

	"git/ymoldabe/forum/models"
)

type Authorization interface {
	InsertUser(form *models.UserSignUp) error
	Authenticate(form *models.UserSignIn) (int, error)
	UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error
	DeleteToken(sessionToken string) error
	GetIdInSessions(sessionToken string) (int, error)
	CheckSessions(userId int) (bool, error)
	UpdateToken(sessionToken string, user_id int) error
	GetTokenSession(cookieToken string) (bool, error)
	GoogleAuthUser(form *models.GoogleLoginUserData) (int, error)
}
