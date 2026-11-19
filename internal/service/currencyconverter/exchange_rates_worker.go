package currencyconverter

import (
	"context"
	"currency-converter/internal/models"
	"slices"
	"sync"
	"time"

	"go.uber.org/zap"
)

func (s *converterService) Start(ctx context.Context) {
	ticker := time.NewTicker(s.workerCfg.WorkerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processExchangeRates(ctx)
		}
	}
}

func (s *converterService) processExchangeRates(ctx context.Context) {
	s.logger.Info("job started sir")

	from := time.Now().Add(-s.workerCfg.WorkerRequestDate).Format(time.DateOnly)
	to := time.Now().Format(time.DateOnly)

	if len(s.workerCfg.WorkerBaseCurrencies) == 0 {
		s.logger.Error("empty base currencies not allowed")
		return
	}

	var baseCurrenciesCh = make(chan string)
	go func() {
		defer close(baseCurrenciesCh)
		for _, currency := range s.workerCfg.WorkerBaseCurrencies {
			select {
			case <-ctx.Done():
				return
			case baseCurrenciesCh <- currency:
			}
		}
	}()

	wg := sync.WaitGroup{}

	for id := 1; id <= s.workerCfg.WorkerPoolSize; id++ {
		wg.Go(func() {
			select {
			case <-ctx.Done():
				return
			case baseCurrency, ok := <-baseCurrenciesCh:
				if !ok {
					return
				}

				currencyRates, err := s.currencyRates.GetCurrencyRatesByDate(ctx, from, to, baseCurrency)
				if err != nil {
					s.logger.Error("failed to fetch currency rates by date",
						zap.Error(err),
						zap.Int("workerID", id),
						zap.String("from", from),
						zap.String("to", to),
						zap.String("baseCurrency", baseCurrency))
					return
				}

				if len(currencyRates) == 0 {
					s.logger.Warn("currency rates not found", zap.Int("workerID", id))
					return
				}

				for chunk := range slices.Chunk(currencyRates, s.workerCfg.WorkerBatchSize) {
					exchangeRates := make([]models.ExchangeRate, 0, len(chunk))
					for _, v := range chunk {
						if v.Base == v.Quote {
							continue
						}

						exchangeRates = append(exchangeRates, models.ExchangeRate{
							CurrencyFrom: v.Base,
							CurrencyTo:   v.Quote,
							Rate:         v.Rate,
							Date:         v.Date,
						})
					}

					if err := s.exchangeRates.CreateExchangeRates(ctx, exchangeRates); err != nil {
						s.logger.Error("failed to create exchange rates",
							zap.Error(err),
							zap.Int("workerID", id),
							zap.Any("exchangeRates", exchangeRates))
						continue
					}
				}
			}
		})
	}

	wg.Wait()
	s.logger.Info("job completed sir")
}

func (s *converterService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}
