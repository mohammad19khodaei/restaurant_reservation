package user

import "context"

type Repository interface {
	Register(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
}
