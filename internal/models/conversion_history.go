package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ConversionHistory struct {
	ID           int
	UserID       int
	CurrencyFrom string
	CurrencyTo   string
	Amount       decimal.Decimal
	Result       decimal.Decimal
	Rate         decimal.Decimal
	CreatedAt    time.Time
}
