package service

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/service/auth"
)

type AuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, loginRequest *auth.LoginRequest) (auth.TokenPairs, error)
}
