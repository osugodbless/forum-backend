package restapi

import (
	"context"

	"github.com/osugodbless/forum-backend/internal/auth/domain"
)

type authService interface {
	SignUp(ctx context.Context, input domain.SignupInput) (*domain.User, error)
}
