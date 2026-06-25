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
