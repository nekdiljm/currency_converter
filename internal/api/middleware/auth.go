package middleware

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/service"
	"net/http"
	"strings"

	"currency-converter/internal/api/response"
	"currency-converter/pkg/logger"
)

type claimsKey struct{}

type MW struct {
	auth   service.AuthService
	logger logger.Logger
}

func New(auth service.AuthService, logger logger.Logger) *MW {
	return &MW{
		auth:   auth,
		logger: logger,
	}
}

func (m *MW) authorize() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				response.Error(w, models.ErrUnauthorized)
				return
			}

			userID, err := m.auth.Authenticate(r.Context(), token)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *MW) ApplyAuth(handlerFunc http.HandlerFunc) http.Handler {
	return m.authorize()(handlerFunc)
}

func UserIDFromCtx(ctx context.Context) int {
	userID, _ := ctx.Value(claimsKey{}).(int)
	return userID
}
