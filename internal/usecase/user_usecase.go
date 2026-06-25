package usecase

import (
	"context"
	"errors"

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
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
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
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.FindByID(ctx, id)
}
