package models

import (
	"fmt"
	"strings"

	"golang.org/x/text/currency"
)

type FavoritePair struct {
	ID           int
	UserID       int
	CurrencyFrom string
	CurrencyTo   string
}

func (f *FavoritePair) Validate() error {
	if f.CurrencyFrom == "" || f.CurrencyTo == "" {
		return ErrInvalidCurrency
	}

	_, err := currency.ParseISO(strings.ToUpper(f.CurrencyFrom))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", f.CurrencyFrom, err)
	}

	_, err = currency.ParseISO(strings.ToUpper(f.CurrencyTo))
	if err != nil {
		return fmt.Errorf("invalid currency code: %q: %w", f.CurrencyTo, err)
	}

	return nil
}
