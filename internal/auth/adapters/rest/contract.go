package rest

import (
	"context"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
)

type AuthService interface {
	SignUp(ctx context.Context, input domain.SignupInput) (*domain.User, error)
}
