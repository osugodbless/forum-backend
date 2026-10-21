package auth

import (
	"context"
	"log/slog"
)

type Service struct {
	userModel *UserModel
	logger    *slog.Logger
}

func NewService(logger *slog.Logger, usermodel *UserModel) *Service {
	return &Service{
		logger:    logger,
		userModel: usermodel,
	}
}

func (s *Service) SignUp(ctx context.Context, dto SignupDTO) (*User, error) {
	return nil, nil
}

func (s *Service) Login() {

}
