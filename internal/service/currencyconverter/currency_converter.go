package currencyconverter

import (
	"context"
	"currency-converter/internal/client/currencyratesprovider"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"currency-converter/internal/service"
	"currency-converter/pkg/logger"
	"currency-converter/pkg/redis"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type converterService struct {
	currencyRates     currencyratesprovider.CurrencyRatesProvider
	conversionHistory repository.ConversionHistory
	cache             redis.Cache
	logger            logger.Logger
}

func NewService(
	currencyRates currencyratesprovider.CurrencyRatesProvider,
	conversionHistory repository.ConversionHistory,
	cache redis.Cache,
	logger logger.Logger,
) service.CurrencyConverter {
	return &converterService{
		currencyRates:     currencyRates,
		conversionHistory: conversionHistory,
		cache:             cache,
		logger:            logger,
	}
}

func (s *converterService) Convert(ctx context.Context, request models.ConvertRequest) (models.ConvertResponse, error) {
	if err := request.Validate(); err != nil {
		s.logger.Warn("failed to validate convert request", zap.Error(err))
		return models.ConvertResponse{}, err
	}

	var currencyRate currencyratesprovider.CurrencyRate
	key := fmt.Sprintf("currency-rate:%s-%s", request.CurrencyFrom, request.CurrencyTo)
	if err := s.cache.Get(ctx, key, &currencyRate); err != nil {
		currencyRate, err = s.currencyRates.GetCurrencyRatesPair(ctx, request.CurrencyFrom, request.CurrencyTo)
		if err != nil {
			s.logger.Error("failed to fetch currency pair", zap.Error(err), zap.Any("request", request))
			return models.ConvertResponse{}, fmt.Errorf("failed to fetch currency pair: %w", err)
		}

		if err = s.cache.Set(ctx, key, &currencyRate, time.Minute); err != nil {
			s.logger.Error("failed to save currency rate in cache", zap.Error(err))
		}
	}

	convertedAmount := request.Amount.Mul(currencyRate.Rate)

	history := models.ConversionHistory{
		UserID:       request.UserID,
		CurrencyFrom: currencyRate.Base,
		CurrencyTo:   currencyRate.Quote,
		Amount:       request.Amount,
		Result:       convertedAmount,
		Rate:         currencyRate.Rate,
	}

	if err := s.conversionHistory.CreateConversionHistory(ctx, history); err != nil {
		s.logger.Error("failed to create conversion history", zap.Error(err), zap.Any("history", history))
		return models.ConvertResponse{}, fmt.Errorf("failed to create conversion history: %w", err)
	}

	return models.ConvertResponse{
		ConvertedAmount: convertedAmount,
		Rate:            currencyRate.Rate,
		From:            currencyRate.Base,
		To:              currencyRate.Quote,
	}, nil
}

func (s *converterService) GetHistoryByUserID(ctx context.Context, userID, limit, cursor int) ([]models.ConversionHistory, error) {
	if limit <= 0 {
		limit = 10
	}
	if cursor < 0 {
		cursor = 0
	}

	history, err := s.conversionHistory.GetConversionHistoryByUserID(ctx, userID, limit, cursor)
	if err != nil {
		if errors.Is(err, models.ErrHistoryNotFound) {
			s.logger.Warn("conversion history not found")
			return nil, err
		}
		s.logger.Error("failed to get conversion history", zap.Error(err))
		return nil, fmt.Errorf("failed to get conversion history: %w", err)
	}

	return history, nil
}

func (s *converterService) GetCurrencyRatesByDate(
	ctx context.Context,
	request models.CurrencyRatesByDateRequest,
) ([]models.CurrencyRatesByDateResponse, error) {
	if err := request.Validate(); err != nil {
		s.logger.Warn("failed to validate currency rates by date", zap.Error(err), zap.Any("request", request))
		return nil, err
	}

	var currencyRatesByDate []currencyratesprovider.CurrencyRate
	key := fmt.Sprintf("currency-rates:%s-%s-%s", request.FromDate, request.ToDate, request.BaseCurrency)

	if err := s.cache.Get(ctx, key, &currencyRatesByDate); err != nil {
		currencyRatesByDate, err = s.currencyRates.GetCurrencyRatesByDate(ctx, request.FromDate, request.ToDate, request.BaseCurrency)
		if err != nil {
			s.logger.Error("failed to get currency rates by date", zap.Error(err), zap.Any("request", request))
			return nil, err
		}

		if err := s.cache.Set(ctx, key, &currencyRatesByDate, time.Minute); err != nil {
			s.logger.Error("failed to save currency rates in cache", zap.Error(err))
		}
	}

	currencyRatesResp := make([]models.CurrencyRatesByDateResponse, 0, len(currencyRatesByDate))
	for _, v := range currencyRatesByDate {
		currencyRatesResp = append(currencyRatesResp, models.CurrencyRatesByDateResponse{
			Date:         v.Date,
			CurrencyFrom: v.Base,
			CurrencyTo:   v.Quote,
			Rate:         v.Rate,
		})
	}

	return currencyRatesResp, nil
}
