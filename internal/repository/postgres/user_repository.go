package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/itskyana/belajar-grpc-go/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`

	err := r.pool.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_email_key") {
				return domain.ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
