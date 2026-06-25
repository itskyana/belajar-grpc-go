package postgres

import (
	"context"
	"errors"

	"github.com/itskyana/belajar-grpc-go/internal/domain"
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
		return errors.New("failed to create user: " + err.Error())
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
