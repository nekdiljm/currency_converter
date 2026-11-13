package conversionhistory

import (
	"context"
	"currency-converter/internal/models"
	"currency-converter/internal/repository"
	"currency-converter/pkg/postgres"
)

type convHistoryRepository struct {
	db postgres.DbCommands
}

func NewRepository(db postgres.DbCommands) repository.ConversionHistory {
	return &convHistoryRepository{db: db}
}

func (c *convHistoryRepository) CreateConversionHistory(
	ctx context.Context,
	history models.ConversionHistory,
) error {
	_, err := c.db.Exec(ctx,
		`INSERT INTO conversion_history (user_id, currency_from, currency_to, amount, result, rate)
			 VALUES($1, $2, $3, $4, $5, $6)`,
		history.UserID,
		history.CurrencyFrom,
		history.CurrencyTo,
		history.Amount,
		history.Result,
		history.Rate)
	if err != nil {
		return err
	}

	return nil
}

func (c *convHistoryRepository) GetConversionHistoryByUserID(
	ctx context.Context,
	userID,
	limit,
	cursor int,
) ([]models.ConversionHistory, error) {
	rows, err := c.db.Query(ctx,
		`SELECT 
				id, 
				user_id, 
				currency_from,
				currency_to,
				amount, 
				result,
				rate, 
				created_at FROM conversion_history
			WHERE user_id = $1 AND (id < $2 OR $2 = 0 )
			ORDER BY created_at DESC
			LIMIT $3
   `, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]models.ConversionHistory, 0, limit)

	for rows.Next() {
		var item models.ConversionHistory
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.CurrencyFrom,
			&item.CurrencyTo,
			&item.Amount,
			&item.Result,
			&item.Rate,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		history = append(history, item)
	}

	if len(history) == 0 {
		return nil, models.ErrHistoryNotFound
	}

	return history, nil
}
