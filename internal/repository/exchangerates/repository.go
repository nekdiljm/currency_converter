package exchangerates

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"currency-converter/pkg/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type exchangeRatesRepository struct {
	db postgres.DbCommands
}

func NewRepository(db postgres.DbCommands) repository.ExchangeRates {
	return &exchangeRatesRepository{
		db: db,
	}
}

func (e *exchangeRatesRepository) CreateExchangeRates(ctx context.Context, rates []models.ExchangeRate) error {
	batch := &pgx.Batch{}

	for _, v := range rates {
		batch.Queue(
			`INSERT INTO exchange_rates (currency_from, currency_to, rate, date)
			 VALUES ($1, $2, $3, $4) 
			 ON CONFLICT(currency_from, currency_to, date) DO UPDATE SET
			 rate = $3, updated_at = now()`, v.CurrencyFrom, v.CurrencyTo, v.Rate, v.Date)
	}

	results := e.db.SendBatch(ctx, batch)
	defer func() {
		_ = results.Close()
	}()

	for range len(rates) {
		_, err := results.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *exchangeRatesRepository) GetExchangeRates(ctx context.Context, currencyFrom, currencyTo string) (models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate

	err := e.db.QueryRow(ctx,
		`SELECT currency_from, currency_to, rate, date FROM exchange_rates
			 WHERE currency_from = $1 AND currency_to = $2
			 ORDER BY date DESC, updated_at DESC LIMIT 1`, currencyFrom, currencyTo).Scan(
		&exchangeRate.CurrencyFrom,
		&exchangeRate.CurrencyTo,
		&exchangeRate.Rate,
		&exchangeRate.Date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ExchangeRate{}, models.ErrRateNotFound
		}

		return models.ExchangeRate{}, fmt.Errorf("failed to get exchange rates from db: %w", err)
	}

	return exchangeRate, nil
}
