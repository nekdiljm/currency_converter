package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ConversionHistory struct {
	ID             int
	UserID         int
	CurrencyFromID int
	CurrencyToID   int
	Amount         decimal.Decimal
	Result         decimal.Decimal
	CreatedAt      time.Time
}
