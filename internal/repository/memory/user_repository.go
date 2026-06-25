package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/itskyana/belajar-grpc-go/internal/domain"
)

type UserRepository struct {
	mu     sync.Mutex
	users  map[int64]*domain.User
	nextID int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[user.ID] = user
	r.nextID++

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}
