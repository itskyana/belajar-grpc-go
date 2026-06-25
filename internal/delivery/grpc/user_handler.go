package grpc

import (
	"context"

	"github.com/itskyana/belajar-grpc-go/gen/userpb"
	"github.com/itskyana/belajar-grpc-go/internal/usecase"
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
		return nil, err
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
		return nil, err
	}

	return &userpb.GetUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
