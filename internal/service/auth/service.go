package auth

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/service"
	"errors"
	"fmt"
	"time"
	"uuid"

	"go.uber.org/zap"

	"currency-converter/internal/config"
	"currency-converter/internal/repository"
	"currency-converter/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type authService struct {
	userRepo   repository.UserRepository
	session    repository.SessionsRepository
	logger     logger.Logger
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(
	cfg *config.Config,
	repo repository.UserRepository,
	session repository.SessionsRepository,
	logger logger.Logger,
) service.AuthService {
	return &authService{
		userRepo:   repo,
		session:    session,
		logger:     logger,
		jwtSecret:  []byte(cfg.JwtConfig.JwtSecret),
		accessTTL:  cfg.JwtConfig.JwtAccessTTL,
		refreshTTL: cfg.JwtConfig.JwtRefreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, registerRequest models.RegisterRequest) (models.User, error) {
	if err := registerRequest.Validate(); err != nil {
		s.logger.Warn("error on request validation", zap.Error(err), zap.Any("request", registerRequest))
		return models.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to generate password", zap.Error(err))
		return models.User{}, fmt.Errorf("generate password hash: %w", err)
	}

	user := models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  string(passwordHash),
	}

	user.ID, err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return models.User{}, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, loginRequest models.LoginRequest) (models.TokenPairs, error) {
	if err := loginRequest.Validate(); err != nil {
		s.logger.Warn("error on request validation", zap.Error(err), zap.Any("request", loginRequest))
		return models.TokenPairs{}, err
	}

	user, err := s.userRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.logger.Warn("user not found by email", zap.Error(err), zap.String("email", loginRequest.Email))
		} else {
			s.logger.Error("failed to get user by email", zap.Error(err))
		}

		return models.TokenPairs{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		s.logger.Warn("error on compare password", zap.Error(err))
		return models.TokenPairs{}, models.ErrInvalidPass
	}

	sessionID := uuid.NewV7().String()
	if err := s.session.Create(ctx, sessionID, user.ID, s.refreshTTL); err != nil {
		s.logger.Error("failed to create session", zap.Error(err), zap.Int("userID", user.ID))
		return models.TokenPairs{}, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := s.generateToken(user.ID, sessionID, accessTokenType, s.accessTTL)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return models.TokenPairs{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user.ID, sessionID, refreshTokenType, s.refreshTTL)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return models.TokenPairs{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return models.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		s.logger.Warn("invalid refresh token")
		return models.ErrInvalidToken
	}

	_, sessionID, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		s.logger.Warn("invalid token", zap.Error(err))
		return err
	}

	if err := s.session.Delete(ctx, sessionID); err != nil && !errors.Is(err, models.ErrSessionNotFound) {
		s.logger.Error("failed to delete session by id", zap.Error(err))
		return err
	}

	return nil
}

func (s *authService) Authenticate(ctx context.Context, accessToken string) (int, error) {
	userID, sessionID, err := s.parse(accessToken, accessTokenType)
	if err != nil {
		s.logger.Warn("invalid token", zap.Error(err))
		return 0, err
	}

	if err := s.session.Exists(ctx, sessionID); err != nil {
		if errors.Is(err, models.ErrSessionNotFound) {
			s.logger.Warn("session not found")
		} else {
			s.logger.Error("failed to check sessionID", zap.Error(err))
		}
		return 0, err
	}

	return userID, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (models.TokenPairs, error) {
	userID, sessionID, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		s.logger.Warn("invalid token", zap.Error(err))
		return models.TokenPairs{}, err
	}

	if err := s.session.Delete(ctx, sessionID); err != nil {
		if errors.Is(err, models.ErrSessionNotFound) {
			s.logger.Warn("session not found")
			return models.TokenPairs{}, models.ErrInvalidToken
		}

		s.logger.Error("failed to delete session by id", zap.Error(err))
		return models.TokenPairs{}, err
	}

	sessionID = uuid.NewV7().String()
	if err := s.session.Create(ctx, sessionID, userID, s.refreshTTL); err != nil {
		s.logger.Error("failed to create session", zap.Error(err), zap.Int("userID", userID))
		return models.TokenPairs{}, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := s.generateToken(userID, sessionID, accessTokenType, s.accessTTL)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return models.TokenPairs{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = s.generateToken(userID, sessionID, refreshTokenType, s.refreshTTL)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return models.TokenPairs{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return models.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
