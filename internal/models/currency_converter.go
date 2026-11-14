package models

import (
	"fmt"
	"strings"
	"time"

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

type CurrencyRatesByDateRequest struct {
	FromDate     string
	ToDate       string
	BaseCurrency string
}

func (c *CurrencyRatesByDateRequest) Validate() error {
	fromDate, err := time.Parse(time.DateOnly, c.FromDate)
	if err != nil {
		return ErrInvalidDateFormat
	}

	toDate, err := time.Parse(time.DateOnly, c.ToDate)
	if err != nil {
		return ErrInvalidDateFormat
	}

	if fromDate.After(toDate) {
		return ErrInvalidDateFormat
	}

	if c.BaseCurrency == "" {
		return ErrInvalidCurrency
	}

	_, err = currency.ParseISO(strings.ToUpper(c.BaseCurrency))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", c.BaseCurrency, err)
	}

	c.FromDate = fromDate.Format(time.DateOnly)
	c.ToDate = toDate.Format(time.DateOnly)

	return nil
}

type CurrencyRatesByDateResponse struct {
	Date         string
	CurrencyFrom string
	CurrencyTo   string
	Rate         decimal.Decimal
}
