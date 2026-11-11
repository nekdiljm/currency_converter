package session

import (
	"context"
	"fmt"
	"time"

	"currency-converter/internal/models"
	"currency-converter/pkg/redis"
)

type Repository struct {
	cache redis.Cache
}

func NewRepository(redis redis.Cache) *Repository {
	return &Repository{cache: redis}
}

func key(sessionID string) string {
	return "session:" + sessionID
}

func (r *Repository) Create(ctx context.Context, sessionID string, userID int, ttl time.Duration) error {
	if err := r.cache.Set(ctx, key(sessionID), userID, ttl); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *Repository) Exists(ctx context.Context, sessionID string) error {
	exists, err := r.cache.Exists(ctx, key(sessionID))
	if err != nil {
		return fmt.Errorf("failed to fetch session: %w", err)
	}

	if !exists {
		return models.ErrSessionNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, sessionID string) error {
	if err := r.cache.Delete(ctx, key(sessionID)); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
