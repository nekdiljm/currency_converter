package models

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	maxRuneInName = 100
	maxRuneInPass = 300
)

type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (r *RegisterRequest) Validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.FirstName == "" || r.LastName == "" || r.Email == "" || r.Password == "" {
		return ErrInvalidInput
	}

	if utf8.RuneCountInString(r.FirstName) > maxRuneInName {
		return ErrFirstNameTooLong
	}

	if addr, err := mail.ParseAddress(r.Email); err != nil || addr.Address != r.Email {
		return ErrInvalidEmail
	}

	if utf8.RuneCountInString(r.LastName) > maxRuneInName {
		return ErrLastNameTooLong
	}

	if utf8.RuneCountInString(r.Password) > maxRuneInPass {
		return ErrPasswordTooLong
	}
	return nil
}

type TokenPairs struct {
	AccessToken  string
	RefreshToken string
}
type LoginRequest struct {
	Email    string
	Password string
}

func (l *LoginRequest) Validate() error {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
	if l.Email == "" || l.Password == "" {
		return ErrInvalidCredentials
	}

	if utf8.RuneCountInString(l.Password) > maxRuneInPass {
		return ErrPasswordTooLong
	}

	return nil
}
