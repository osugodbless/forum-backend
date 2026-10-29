package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store  AuthStore
	logger *slog.Logger
}

func NewAuthService(store AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

type SignupInput struct {
	FullName string
	Email    string
	Password string
}

func (s *AuthService) SignUp(ctx context.Context, dto SignupInput) (*User, error) {
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

func (s *AuthService) Login() {

}
