package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	ID             int
	CurrencyFromID int
	CurrencyToID   int
	ConvRate       decimal.Decimal
	UpdatedAt      time.Time
}
