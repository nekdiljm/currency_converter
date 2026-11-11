package user

import (
	"context"
	models2 "currency-converter/internal/repository/models"
	"errors"
	"fmt"

	"currency-converter/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolationCode = "23505"

type Repository struct {
	db postgres.DbCommands
}

func NewRepository(db postgres.DbCommands) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *models2.User) (int, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (first_name, last_name, email, password)
			 VALUES ($1, $2, $3,$4)
			 RETURNING id, created_at;`, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return 0, models2.ErrEmailAlreadyExists
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return user.ID, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models2.User, error) {
	var user models2.User
	err := r.db.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, password, created_at FROM users
			 WHERE email=$1`, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models2.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}
