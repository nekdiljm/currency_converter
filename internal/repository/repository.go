package repository

import (
	"context"
	"currency-converter/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
