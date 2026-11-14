package v1

import (
	"currency-converter/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

type registerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r registerRequest) toInput() models.RegisterRequest {
	return models.RegisterRequest{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  r.Password,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l loginRequest) toInput() models.LoginRequest {
	return models.LoginRequest{
		Email:    l.Email,
		Password: l.Password,
	}
}

type registerResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func toUserResponse(in models.User) registerResponse {
	return registerResponse{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		CreatedAt: in.CreatedAt.Format(time.RFC3339),
	}
}

type tokenPairsResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func toTokenPairsResponse(in models.TokenPairs) tokenPairsResponse {
	return tokenPairsResponse{
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
	}
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type convertRequest struct {
	Amount       decimal.Decimal `json:"amount"`
	CurrencyFrom string          `json:"currencyFrom"`
	CurrencyTo   string          `json:"currencyTo"`
}

func (c *convertRequest) toInput() models.ConvertRequest {
	return models.ConvertRequest{
		Amount:       c.Amount,
		CurrencyFrom: c.CurrencyFrom,
		CurrencyTo:   c.CurrencyTo,
	}
}

type convertResponse struct {
	ConvertedAmount decimal.Decimal `json:"convertedAmount"`
	Rate            decimal.Decimal `json:"rate"`
	CurrencyFrom    string          `json:"currencyFrom"`
	CurrencyTo      string          `json:"currencyTo"`
}

func toConvertResponse(in models.ConvertResponse) convertResponse {
	return convertResponse{
		ConvertedAmount: in.ConvertedAmount,
		Rate:            in.Rate,
		CurrencyFrom:    in.From,
		CurrencyTo:      in.To,
	}
}

type historyResponse struct {
	ID           int             `json:"id"`
	CurrencyFrom string          `json:"currencyFrom"`
	CurrencyTo   string          `json:"currencyTo"`
	Amount       decimal.Decimal `json:"amount"`
	Result       decimal.Decimal `json:"result"`
	Rate         decimal.Decimal `json:"rate"`
	CreatedAt    string          `json:"createdAt"`
}

func toHistoryResponse(in []models.ConversionHistory) []historyResponse {
	out := make([]historyResponse, 0, len(in))

	for _, v := range in {
		out = append(out, historyResponse{
			ID:           v.ID,
			CurrencyFrom: v.CurrencyFrom,
			CurrencyTo:   v.CurrencyTo,
			Amount:       v.Amount,
			Result:       v.Result,
			Rate:         v.Rate,
			CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		})
	}

	return out
}

type favoritePairCreateRequest struct {
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
}

func (f *favoritePairCreateRequest) toInput() models.FavoritePair {
	return models.FavoritePair{
		CurrencyFrom: f.CurrencyFrom,
		CurrencyTo:   f.CurrencyTo,
	}
}

type favoritePairCreateResponse struct {
	ID int `json:"id"`
}

func toFavoritePairCreateResponse(id int) favoritePairCreateResponse {
	return favoritePairCreateResponse{
		ID: id,
	}
}

type favoritePairsResponse struct {
	ID           int    `json:"id"`
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
}

func toFavoritePairsResponse(in []models.FavoritePair) []favoritePairsResponse {
	out := make([]favoritePairsResponse, 0, len(in))

	for _, v := range in {
		out = append(out, favoritePairsResponse{
			ID:           v.ID,
			CurrencyFrom: v.CurrencyFrom,
			CurrencyTo:   v.CurrencyTo,
		})
	}

	return out
}

type favoritePairsDeleteRequest struct {
	ID int `json:"id"`
}

type currencyRatesByDateResponse struct {
	Date         string          `json:"date"`
	CurrencyFrom string          `json:"currencyFrom"`
	CurrencyTo   string          `json:"currencyTo"`
	Rate         decimal.Decimal `json:"rate"`
}

func toCurrencyRatesByDateResponse(in []models.CurrencyRatesByDateResponse) []currencyRatesByDateResponse {
	out := make([]currencyRatesByDateResponse, 0, len(in))

	for _, v := range in {
		out = append(out, currencyRatesByDateResponse{
			Date:         v.Date,
			CurrencyFrom: v.CurrencyFrom,
			CurrencyTo:   v.CurrencyTo,
			Rate:         v.Rate,
		})
	}

	return out
}
