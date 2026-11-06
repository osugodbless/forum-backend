package profile

import (
	"context"
)

type profileServiceManager interface {
}

type profileStore interface {
	InsertDefaultProfile(ctx context.Context, profile *UserProfile) error
}
