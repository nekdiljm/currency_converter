package v1

import (
	"currency-converter/internal/service"
	"currency-converter/pkg/logger"
)

type Handler struct {
	authService       service.AuthService
	currencyConverter service.CurrencyConverter
	favoritePairs     service.FavoritePairs
	logger            logger.Logger
}

func NewHandler(
	authService service.AuthService,
	currencyConverter service.CurrencyConverter,
	favoritePairs service.FavoritePairs,
	logger logger.Logger,
) *Handler {
	return &Handler{
		authService:       authService,
		currencyConverter: currencyConverter,
		favoritePairs:     favoritePairs,
		logger:            logger,
	}
}
