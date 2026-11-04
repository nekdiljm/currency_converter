package models

import "errors"

var (
	ErrInvalidInput = errors.New("error invalid input")
	ErrInvalidPass  = errors.New("error invalid password")
	ErrUserNotFound = errors.New("error user not found")
)
