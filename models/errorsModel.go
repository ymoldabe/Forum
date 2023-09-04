package models

import "errors"

var (
	ErrFormNotValid       = errors.New("Form not valid")
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("UNIQUE constraint failed: users.email")
	ErrNoRowsInResSet     = errors.New("no rows in result set")
	ErrDeleteFailed       = errors.New("delete failed")
)
