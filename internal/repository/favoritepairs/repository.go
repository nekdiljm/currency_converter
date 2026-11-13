package favoritepairs

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"currency-converter/pkg/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

type favoritePairsRepository struct {
	db postgres.DbCommands
}

func NewRepository(db postgres.DbCommands) repository.FavoritePairs {
	return &favoritePairsRepository{db: db}
}

func (f *favoritePairsRepository) CreateFavoritePair(ctx context.Context, favPair models.FavoritePair) (int, error) {
	err := f.db.QueryRow(ctx,
		`INSERT INTO favorite_pairs (user_id, currency_from, currency_to) 
			 VALUES ($1, $2, $3)
			 RETURNING id`, favPair.UserID, favPair.CurrencyFrom, favPair.CurrencyTo).Scan(&favPair.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return 0, models.ErrFavoritePairAlreadyExists
		}
		return 0, fmt.Errorf("add favorite pair: %w", err)
	}

	return favPair.ID, nil
}

func (f *favoritePairsRepository) GetFavoritePairsByUserID(
	ctx context.Context, userID, limit, offset int,
) ([]models.FavoritePair, error) {
	rows, err := f.db.Query(ctx,
		`SELECT id, user_id, currency_from, currency_to FROM favorite_pairs
			 WHERE user_id = $1
			 LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favoritePairs := make([]models.FavoritePair, 0, limit)

	for rows.Next() {
		var item models.FavoritePair
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.CurrencyFrom,
			&item.CurrencyTo,
		); err != nil {
			return nil, err
		}

		favoritePairs = append(favoritePairs, item)
	}

	if len(favoritePairs) == 0 {
		return nil, models.ErrFavoritePairNotFound
	}

	return favoritePairs, nil
}

func (f *favoritePairsRepository) DeleteFavoritePairByID(ctx context.Context, id int) error {
	result, err := f.db.Exec(ctx, `DELETE FROM favorite_pairs WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrFavoritePairNotFound
	}

	return nil
}
