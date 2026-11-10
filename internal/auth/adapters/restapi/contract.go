package restapi

import (
	"context"

	"github.com/osugodbless/forum-backend/internal/auth/domain"
)

type authService interface {
	SignUp(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
}
