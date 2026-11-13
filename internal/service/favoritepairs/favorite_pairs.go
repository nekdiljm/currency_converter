package favoritepairs

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"currency-converter/internal/service"
	"currency-converter/pkg/logger"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type favoritePairsService struct {
	favoritePairs repository.FavoritePairs
	logger        logger.Logger
}

func NewService(favoritePairs repository.FavoritePairs, logger logger.Logger) service.FavoritePairs {
	return &favoritePairsService{
		favoritePairs: favoritePairs,
		logger:        logger,
	}
}

func (f *favoritePairsService) CreateFavPair(ctx context.Context, pair models.FavoritePair) (int, error) {
	if err := pair.Validate(); err != nil {
		f.logger.Warn("failed to validate favorite pair", zap.Error(err))
		return 0, err
	}

	favPairID, err := f.favoritePairs.CreateFavoritePair(ctx, pair)
	if err != nil {
		f.logger.Error("failed to create favorite pair", zap.Error(err), zap.Any("pair", pair))
		return 0, err
	}

	return favPairID, nil
}

func (f *favoritePairsService) GetFavoritePairByUserID(ctx context.Context, userID, limit, offset int) ([]models.FavoritePair, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	favoritePairs, err := f.favoritePairs.GetFavoritePairsByUserID(ctx, userID, limit, offset)
	if err != nil {
		if errors.Is(err, models.ErrFavoritePairNotFound) {
			f.logger.Warn("conversion history not found")
			return nil, err
		}
		f.logger.Error("failed to get favorite pairs", zap.Error(err))
		return nil, fmt.Errorf("failed to get favorite pairs: %w", err)
	}

	return favoritePairs, nil
}

func (f *favoritePairsService) DeleteFavoritePair(ctx context.Context, id int) error {
	if err := f.favoritePairs.DeleteFavoritePairByID(ctx, id); err != nil {
		if errors.Is(err, models.ErrFavoritePairNotFound) {
			f.logger.Warn("favorite pair not found", zap.Error(err), zap.Int("id", id))
			return err
		}

		f.logger.Error("failed to delete favorite pairs", zap.Error(err))
		return fmt.Errorf("failed to delete favorite pairs: %w", err)
	}

	return nil
}
