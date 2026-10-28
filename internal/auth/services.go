package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store  *UserModel
	logger *slog.Logger
}

func NewService(logger *slog.Logger, usermodel *UserModel) *Service {
	return &Service{
		logger: logger,
		store:  usermodel,
	}
}

func (s *Service) SignUp(ctx context.Context, dto SignupDTO) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &User{
		ID:           uuid.New(),
		Email:        dto.Email,
		PasswordHash: string(bytes),
	}

	err = s.store.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user service: %w", err)
	}

	return user, nil
}

func (s *Service) Login() {

}
