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

func toUserResponse(u models.User) registerResponse {
	return registerResponse{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

type tokenPairsResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func toTokenPairsResponse(tokens models.TokenPairs) tokenPairsResponse {
	return tokenPairsResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
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

func toConvertResponse(response models.ConvertResponse) convertResponse {
	return convertResponse{
		ConvertedAmount: response.ConvertedAmount,
		Rate:            response.Rate,
		CurrencyFrom:    response.From,
		CurrencyTo:      response.To,
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

func toHistoryResponse(history []models.ConversionHistory) []historyResponse {
	historyResp := make([]historyResponse, 0, len(history))

	for _, v := range history {
		historyResp = append(historyResp, historyResponse{
			ID:           v.ID,
			CurrencyFrom: v.CurrencyFrom,
			CurrencyTo:   v.CurrencyTo,
			Amount:       v.Amount,
			Result:       v.Result,
			Rate:         v.Rate,
			CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		})
	}

	return historyResp
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

func toFavoritePairsResponse(pairs []models.FavoritePair) []favoritePairsResponse {
	favoritePairsResp := make([]favoritePairsResponse, 0, len(pairs))

	for _, v := range pairs {
		favoritePairsResp = append(favoritePairsResp, favoritePairsResponse{
			ID:           v.ID,
			CurrencyFrom: v.CurrencyFrom,
			CurrencyTo:   v.CurrencyTo,
		})
	}

	return favoritePairsResp
}

type favoritePairsDeleteRequest struct {
	ID int `json:"id"`
}
