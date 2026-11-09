package auth

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxRuneInName = 100
	maxRuneInPass = 300
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewService(repo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (r *RegisterRequest) validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.FirstName == "" || r.LastName == "" || r.Email == "" || r.Password == "" {
		return models.ErrInvalidInput
	}

	if utf8.RuneCountInString(r.FirstName) > maxRuneInName {
		return models.ErrFirstNameTooLong
	}

	if addr, err := mail.ParseAddress(r.Email); err != nil || addr.Address != r.Email {
		return models.ErrInvalidEmail
	}

	if utf8.RuneCountInString(r.LastName) > maxRuneInName {
		return models.ErrLastNameTooLong
	}

	if utf8.RuneCountInString(r.Password) > maxRuneInPass {
		return models.ErrPasswordTooLong
	}
	return nil
}

func (s *AuthService) Register(ctx context.Context, registerRequest *RegisterRequest) (*models.User, error) {
	if err := registerRequest.validate(); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generate password hash: %w", err)
	}

	user := &models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  string(passwordHash),
	}

	user.ID, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type LoginRequest struct {
	Email    string
	Password string
}

func (l *LoginRequest) validate() error {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
	if l.Email == "" || l.Password == "" {
		return models.ErrInvalidCredentials
	}

	if utf8.RuneCountInString(l.Password) > maxRuneInPass {
		return models.ErrPasswordTooLong
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, loginRequest *LoginRequest) (*models.User, error) {
	if err := loginRequest.validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return nil, models.ErrInvalidPass
	}

	return user, nil
}
