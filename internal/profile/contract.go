package profile

import "context"

type ProfileStore interface {
	InsertProfile(ctx context.Context, profile *UserProfile) error
}
