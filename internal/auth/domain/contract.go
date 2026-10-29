package domain

import "context"

type AuthStore interface {
	Create(ctx context.Context, user *User) error
}

type ProfileCreator interface {
}
