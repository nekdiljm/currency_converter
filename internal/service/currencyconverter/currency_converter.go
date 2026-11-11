package currencyconverter

import (
	"context"
	"currency-converter/internal/client/currencyratesprovider"
	"currency-converter/internal/repository/models"
	"currency-converter/pkg/logger"
	"currency-converter/pkg/redis"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/text/currency"
)

type CurrencyConverter interface {
	Convert(ctx context.Context, request ConvertRequest) (ConvertResponse, error)
}

type service struct {
	currencyRates currencyratesprovider.CurrencyRatesProvider
	cache         redis.Cache
	logger        logger.Logger
}

func NewService(
	currencyRates currencyratesprovider.CurrencyRatesProvider,
	cache redis.Cache,
	logger logger.Logger,
) CurrencyConverter {
	return &service{
		currencyRates: currencyRates,
		cache:         cache,
		logger:        logger,
	}
}

type ConvertRequest struct {
	Amount decimal.Decimal
	From   string
	To     string
}

func (c *ConvertRequest) validate() error {
	if c.Amount.IsZero() || c.Amount.IsNegative() {
		return models.ErrInvalidAmount
	}

	if c.From == "" || c.To == "" {
		return models.ErrInvalidCurrency
	}

	_, err := currency.ParseISO(strings.ToUpper(c.From))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", c.From, err)
	}

	_, err = currency.ParseISO(strings.ToUpper(c.To))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", c.To, err)
	}

	return nil
}

type ConvertResponse struct {
	ConvertedAmount decimal.Decimal
	Rate            decimal.Decimal
	From            string
	To              string
}

func (s *service) Convert(ctx context.Context, request ConvertRequest) (ConvertResponse, error) {
	if err := request.validate(); err != nil {
		s.logger.Warn("failed to validate convert request", zap.Error(err))
		return ConvertResponse{}, err
	}

	currencyRate, err := s.currencyRates.GetCurrencyRatesPair(ctx, request.From, request.To)
	if err != nil {
		s.logger.Error("failed to fetch currency pair", zap.Error(err), zap.Any("request", request))
		return ConvertResponse{}, fmt.Errorf("failed to fetch currency pair: %w", err)
	}

	convertedAmount := request.Amount.Mul(currencyRate.Rate).Truncate(2)

	return ConvertResponse{
		ConvertedAmount: convertedAmount,
		Rate:            currencyRate.Rate,
		From:            currencyRate.Base,
		To:              currencyRate.Quote,
	}, nil
}
