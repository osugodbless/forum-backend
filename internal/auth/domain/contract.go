package domain

import "context"

type authStore interface {
	Create(ctx context.Context, user *User) error
}

type profileCreator interface {
}
