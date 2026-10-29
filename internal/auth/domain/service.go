package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")
)

type AuthService struct {
	store AuthStore
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
