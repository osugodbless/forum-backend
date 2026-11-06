package domain

import (
	"context"

	"github.com/google/uuid"
)

type authStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	UserExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, newHash string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
