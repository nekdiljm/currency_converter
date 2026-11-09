package v1

import (
	"currency-converter/internal/models"
	"currency-converter/internal/service/auth"
	"time"
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
	CreatedAt string `json:"created_at"`
}

func toUserResponse(u *models.User) registerResponse {
	return registerResponse{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func toLoginResponse(tokens auth.TokenPairs) loginResponse {
	return loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
