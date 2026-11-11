package v1

import (
	"currency-converter/internal/service/auth"
	"currency-converter/internal/service/currencyconverter"
	"currency-converter/pkg/logger"
)

type Handler struct {
	authService       auth.AuthService
	currencyConverter currencyconverter.CurrencyConverter
	logger            logger.Logger
}

func NewHandler(
	authService auth.AuthService,
	currencyConverter currencyconverter.CurrencyConverter,
	logger logger.Logger,
) *Handler {
	return &Handler{
		authService:       authService,
		currencyConverter: currencyConverter,
		logger:            logger,
	}
}
