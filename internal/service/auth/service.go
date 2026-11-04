package auth

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxRuneInName = 100
	maxRuneInPass = 300
)

type AuthService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, password string) (*models.User, error) {

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	email = strings.ToLower(strings.TrimSpace(email))

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return nil, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(firstName) > maxRuneInName {
		return nil, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(lastName) > maxRuneInName {
		return nil, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(password) > maxRuneInPass {
		return nil, models.ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generate password hash: %w", err)
	}

	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(passwordHash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	if email == "" || password == "" {
		return nil, models.ErrInvalidInput
	}

	if utf8.RuneCountInString(password) > maxRuneInPass {
		return nil, models.ErrInvalidInput
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, models.ErrInvalidPass
	}

	return user, nil
}
