package usecase

import (
	"context"
	"strings"

	"github.com/itskyana/belajar-grpc-go/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	if email == "" {
		return nil, domain.ErrInvalidInput
	}

	if !strings.Contains(email, "@") {
		return nil, domain.ErrInvalidInput
	}

	user := &domain.User{
		Name:  name,
		Email: email,
	}

	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return u.userRepo.FindByID(ctx, id)
}

func (u *UserUseCase) ListUsers(ctx context.Context, limit, offset int32) ([]*domain.User, int32, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil

}

func (u *UserUseCase) UpdateUser(ctx context.Context, id int64, name, email string) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, domain.ErrInvalidInput
	}

	if !strings.Contains(email, "@") {
		return nil, domain.ErrInvalidInput
	}

	user := &domain.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
