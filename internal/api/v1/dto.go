package v1

import (
	"currency-converter/internal/repository/models"
	"currency-converter/internal/service/currencyconverter"
	"time"

	"currency-converter/internal/service/auth"

	"github.com/shopspring/decimal"
)

type registerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r registerRequest) toInput() *auth.RegisterRequest {
	return &auth.RegisterRequest{
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

func (l loginRequest) toInput() *auth.LoginRequest {
	return &auth.LoginRequest{
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

func toUserResponse(u *models.User) registerResponse {
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

func toTokenPairsResponse(tokens auth.TokenPairs) tokenPairsResponse {
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
	Amount decimal.Decimal `json:"amount"`
	From   string          `json:"from"`
	To     string          `json:"to"`
}

func (c *convertRequest) toInput() currencyconverter.ConvertRequest {
	return currencyconverter.ConvertRequest{
		Amount: c.Amount,
		From:   c.From,
		To:     c.To,
	}
}

type convertResponse struct {
	ConvertedAmount decimal.Decimal `json:"convertedAmount"`
	Rate            decimal.Decimal `json:"rate"`
	From            string          `json:"from"`
	To              string          `json:"to"`
}

func toConvertResponse(response currencyconverter.ConvertResponse) convertResponse {
	return convertResponse{
		ConvertedAmount: response.ConvertedAmount,
		Rate:            response.Rate,
		From:            response.From,
		To:              response.To,
	}
}
