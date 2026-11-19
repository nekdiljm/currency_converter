package session

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"fmt"
	"time"

	"currency-converter/pkg/redis"
)

type sessionRepository struct {
	cache redis.Cache
}

func NewRepository(redis redis.Cache) repository.SessionsRepository {
	return &sessionRepository{cache: redis}
}

func key(sessionID string) string {
	return "session:" + sessionID
}

func (r *sessionRepository) Create(ctx context.Context, sessionID string, userID int, ttl time.Duration) error {
	if err := r.cache.Set(ctx, key(sessionID), userID, ttl); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *sessionRepository) Exists(ctx context.Context, sessionID string) error {
	exists, err := r.cache.Exists(ctx, key(sessionID))
	if err != nil {
		return fmt.Errorf("failed to fetch session: %w", err)
	}

	if !exists {
		return models.ErrSessionNotFound
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	deleted, err := r.cache.Delete(ctx, key(sessionID))
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if deleted == 0 {
		return models.ErrSessionNotFound
	}

	return nil
}
