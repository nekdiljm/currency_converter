package repository

import (
	"context"
	"currency-converter/internal/models"
	"time"
)

type Repository struct {
	UserRepository     UserRepository
	SessionsRepository SessionsRepository
	ConversionHistory  ConversionHistory
	FavoritePairs      FavoritePairs
}

type UserRepository interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, sessionID string, userID int, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) error
	Delete(ctx context.Context, sessionID string) error
}

type ConversionHistory interface {
	CreateConversionHistory(ctx context.Context, history models.ConversionHistory) error
	GetConversionHistoryByUserID(ctx context.Context, userID, limit, cursor int) ([]models.ConversionHistory, error)
}

type FavoritePairs interface {
	CreateFavoritePair(ctx context.Context, favPair models.FavoritePair) (int, error)
	GetFavoritePairsByUserID(ctx context.Context, userID, limit, offset int) ([]models.FavoritePair, error)
	DeleteFavoritePairByID(ctx context.Context, id int) error
}
