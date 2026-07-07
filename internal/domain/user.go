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
	List(ctx context.Context, limit, offset int32) ([]*User, error)
	Count(ctx context.Context) (int32, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
