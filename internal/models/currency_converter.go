package models

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

type ConvertRequest struct {
	UserID       int
	Amount       decimal.Decimal
	CurrencyFrom string
	CurrencyTo   string
}

func (c *ConvertRequest) Validate() error {
	if c.Amount.IsZero() || c.Amount.IsNegative() {
		return ErrInvalidAmount
	}

	if c.CurrencyFrom == "" || c.CurrencyTo == "" {
		return ErrInvalidCurrency
	}

	_, err := currency.ParseISO(strings.ToUpper(c.CurrencyFrom))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", c.CurrencyFrom, err)
	}

	_, err = currency.ParseISO(strings.ToUpper(c.CurrencyTo))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", c.CurrencyTo, err)
	}

	return nil
}

type ConvertResponse struct {
	ConvertedAmount decimal.Decimal
	Rate            decimal.Decimal
	From            string
	To              string
}
