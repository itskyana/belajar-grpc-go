package domain

import "context"

type User struct {
	ID    int64
	Name  string
	Email string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
}
