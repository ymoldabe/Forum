package models

import "git/ymoldabe/forum/validator"

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
