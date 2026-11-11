package currencyratesprovider

import "github.com/shopspring/decimal"

type CurrencyRate struct {
	Date  string          `json:"date"`
	Base  string          `json:"base"`
	Quote string          `json:"quote"`
	Rate  decimal.Decimal `json:"rate"`
}
