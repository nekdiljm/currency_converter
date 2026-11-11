package middleware

import (
	"context"
	"currency-converter/internal/repository/models"
	"currency-converter/internal/service/auth"
	"net/http"
	"strings"

	"currency-converter/internal/api/response"
	"currency-converter/pkg/logger"
)

type claimsKey struct{}

type MW struct {
	auth   auth.AuthService
	logger logger.Logger
}

func New(auth auth.AuthService, logger logger.Logger) *MW {
	return &MW{
		auth:   auth,
		logger: logger,
	}
}

func (m *MW) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				response.Error(w, models.ErrUnauthorized)
				return
			}

			claims, err := m.auth.Authenticate(r.Context(), token)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
