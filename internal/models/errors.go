package models

import "errors"

var (
	ErrInvalidInput              = errors.New("error first name, last name, email and password are required")
	ErrInvalidEmail              = errors.New("error invalid email")
	ErrInvalidPass               = errors.New("error invalid password")
	ErrUserNotFound              = errors.New("error user not found")
	ErrEmailAlreadyExists        = errors.New("error email already exists")
	ErrPasswordTooLong           = errors.New("password must be at most 300 characters")
	ErrFirstNameTooLong          = errors.New("first name must be at most 100 characters")
	ErrLastNameTooLong           = errors.New("last name must be at most 100 characters")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrInvalidToken              = errors.New("invalid token")
	ErrSessionNotFound           = errors.New("session not found")
	ErrInvalidAmount             = errors.New("amount must be greater than zero")
	ErrInvalidCurrency           = errors.New("invalid currency")
	ErrHistoryNotFound           = errors.New("error history not found")
	ErrFavoritePairAlreadyExists = errors.New("favorite pair already exists")
	ErrFavoritePairNotFound      = errors.New("favorite pair not found")
	ErrInvalidDateFormat         = errors.New("invalid date format")
)
