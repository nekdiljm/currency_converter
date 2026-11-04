package auth

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, password string) (*models.User, error) {
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

func ()