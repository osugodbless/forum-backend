package profile

import (
	"context"
	"uuid"
)

type profileServiceManager interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, fullName string) error
}

type profileStore interface {
	InsertProfile(ctx context.Context, profile *UserProfile) error
}
