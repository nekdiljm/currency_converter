package service

import (
	"context"
	"currency-converter/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, request models.RegisterRequest) (models.User, error)
	Login(ctx context.Context, loginRequest models.LoginRequest) (models.TokenPairs, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (models.TokenPairs, error)
	Authenticate(ctx context.Context, accessToken string) (int, error)
}

type CurrencyConverter interface {
	Convert(ctx context.Context, request models.ConvertRequest) (models.ConvertResponse, error)
	GetHistoryByUserID(ctx context.Context, userID, limit, cursor int) ([]models.ConversionHistory, error)
}

type FavoritePairs interface {
	CreateFavPair(ctx context.Context, pair models.FavoritePair) (int, error)
	GetFavoritePairByUserID(ctx context.Context, userID, limit, offset int) ([]models.FavoritePair, error)
	DeleteFavoritePair(ctx context.Context, id int) error
}
