package domain

import "context"

type authStore interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type profileCreator interface {
}
