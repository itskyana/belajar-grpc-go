package grpc

import (
	"context"
	"errors"

	"github.com/itskyana/belajar-grpc-go/gen/userpb"
	"github.com/itskyana/belajar-grpc-go/internal/domain"
	"github.com/itskyana/belajar-grpc-go/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user, err := h.userUseCase.CreateUser(ctx, req.GetName(), req.GetEmail())
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.CreateUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.userUseCase.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.GetUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, "email already exists")
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
