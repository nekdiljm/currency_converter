package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"currency-converter/internal/models"
)

func (s *authService) generateToken(userID int, sessionID string, tokenType string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(userID),
		"jti":  sessionID,
		"type": tokenType,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) parse(token, tokenType string) (int, string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims,
		func(*jwt.Token) (any, error) {
			return s.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return 0, "", fmt.Errorf("%w: %w", models.ErrInvalidToken, err)
	}

	v, ok := claims["type"].(string)
	if !ok || v != tokenType {
		return 0, "", fmt.Errorf("expected %s token, got %q: %w", tokenType, v, models.ErrInvalidToken)

	}

	sessionID, ok := claims["jti"].(string)
	if !ok || sessionID == "" {
		return 0, "", fmt.Errorf("invalid sessionID: %w", models.ErrInvalidToken)
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return 0, "", fmt.Errorf("invalid subject: %w", models.ErrInvalidToken)
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid userID: %w", models.ErrInvalidToken)
	}

	return int(id), sessionID, nil
}
