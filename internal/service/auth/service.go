package auth

import (
	"context"
	models2 "currency-converter/internal/repository/models"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"go.uber.org/zap"

	"currency-converter/internal/config"
	"currency-converter/internal/repository"
	"currency-converter/internal/repository/session"
	"currency-converter/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

const (
	maxRuneInName = 100
	maxRuneInPass = 300
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type AuthService interface {
	Register(ctx context.Context, request *RegisterRequest) (*models2.User, error)
	Login(ctx context.Context, loginRequest *LoginRequest) (TokenPairs, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (TokenPairs, error)
	Authenticate(ctx context.Context, accessToken string) (int, error)
}

type authService struct {
	userRepo   repository.UserRepository
	session    *session.Repository
	logger     logger.Logger
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(
	cfg *config.Config,
	repo repository.UserRepository,
	session *session.Repository,
	logger logger.Logger,
) AuthService {
	return &authService{
		userRepo:   repo,
		session:    session,
		logger:     logger,
		jwtSecret:  []byte(cfg.JwtConfig.JwtSecret),
		accessTTL:  cfg.JwtConfig.JwtAccessTTL,
		refreshTTL: cfg.JwtConfig.JwtRefreshTTL,
	}
}

type RegisterRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (r *RegisterRequest) validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.FirstName == "" || r.LastName == "" || r.Email == "" || r.Password == "" {
		return models2.ErrInvalidInput
	}

	if utf8.RuneCountInString(r.FirstName) > maxRuneInName {
		return models2.ErrFirstNameTooLong
	}

	if addr, err := mail.ParseAddress(r.Email); err != nil || addr.Address != r.Email {
		return models2.ErrInvalidEmail
	}

	if utf8.RuneCountInString(r.LastName) > maxRuneInName {
		return models2.ErrLastNameTooLong
	}

	if utf8.RuneCountInString(r.Password) > maxRuneInPass {
		return models2.ErrPasswordTooLong
	}
	return nil
}

func (s *authService) Register(ctx context.Context, registerRequest *RegisterRequest) (*models2.User, error) {
	if err := registerRequest.validate(); err != nil {
		s.logger.Warn("error on request validation", zap.Error(err), zap.Any("request", registerRequest))
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to generate password", zap.Error(err))
		return nil, fmt.Errorf("generate password hash: %w", err)
	}

	user := &models2.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  string(passwordHash),
	}

	user.ID, err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

type LoginRequest struct {
	Email    string
	Password string
}

type TokenPairs struct {
	AccessToken  string
	RefreshToken string
}

func (l *LoginRequest) validate() error {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
	if l.Email == "" || l.Password == "" {
		return models2.ErrInvalidCredentials
	}

	if utf8.RuneCountInString(l.Password) > maxRuneInPass {
		return models2.ErrPasswordTooLong
	}

	return nil
}

func (s *authService) Login(ctx context.Context, loginRequest *LoginRequest) (TokenPairs, error) {
	if err := loginRequest.validate(); err != nil {
		s.logger.Warn("error on request validation", zap.Error(err), zap.Any("request", loginRequest))
		return TokenPairs{}, err
	}

	user, err := s.userRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, models2.ErrUserNotFound) {
			s.logger.Warn("user not found by email", zap.Error(err), zap.String("email", loginRequest.Email))
		} else {
			s.logger.Error("failed to get user by email", zap.Error(err))
		}

		return TokenPairs{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		s.logger.Warn("error on compare password", zap.Error(err))
		return TokenPairs{}, models2.ErrInvalidPass
	}

	sessionID := uuid.NewV7().String()
	if err := s.session.Create(ctx, sessionID, user.ID, s.refreshTTL); err != nil {
		s.logger.Error("failed to create session", zap.Error(err), zap.Int("userID", user.ID))
		return TokenPairs{}, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := s.generateToken(user.ID, sessionID, accessTokenType, s.accessTTL)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return TokenPairs{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user.ID, sessionID, refreshTokenType, s.refreshTTL)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return TokenPairs{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		s.logger.Warn("invalid refresh token")
		return models2.ErrInvalidToken
	}

	_, sessionID, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		s.logger.Warn("invalid token", zap.Error(err))
		return err
	}

	if err := s.session.Delete(ctx, sessionID); err != nil {
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
		if errors.Is(err, models2.ErrSessionNotFound) {
			s.logger.Warn("session not found")
		} else {
			s.logger.Error("failed to check sessionID", zap.Error(err))
		}
		return 0, err
	}

	return userID, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (TokenPairs, error) {
	userID, sessionID, err := s.parse(refreshToken, refreshTokenType)
	if err != nil {
		s.logger.Warn("invalid token", zap.Error(err))
		return TokenPairs{}, err
	}

	if err := s.session.Delete(ctx, sessionID); err != nil {
		s.logger.Error("failed to delete session by id", zap.Error(err))
		return TokenPairs{}, err
	}

	sessionID = uuid.NewV7().String()
	if err := s.session.Create(ctx, sessionID, userID, s.refreshTTL); err != nil {
		s.logger.Error("failed to create session", zap.Error(err), zap.Int("userID", userID))
		return TokenPairs{}, fmt.Errorf("failed to create session: %w", err)
	}

	accessToken, err := s.generateToken(userID, sessionID, accessTokenType, s.accessTTL)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return TokenPairs{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = s.generateToken(userID, sessionID, refreshTokenType, s.refreshTTL)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return TokenPairs{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
