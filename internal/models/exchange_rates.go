package models

import (
	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	ID           int
	CurrencyFrom string
	CurrencyTo   string
	Rate         decimal.Decimal
	Date         string
}
