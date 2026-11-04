package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/osugodbless/forum-backend.git/internal/platform/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")

	ErrIDGeneration = errors.New("user id generation failed")
)

type AuthService struct {
	store authStore
}

func NewAuthService(store authStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

type SignupInput struct {
	FullName string
	Email    string
	Password string
}

func (s *AuthService) SignUp(ctx context.Context, input SignupInput) (*User, error) {
	// Normalize signup input data
	input.FullName = strings.Join(strings.Fields(input.FullName), " ")
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	v := validator.NewValidator()
	if ValidateUser(v, &input); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrIDGeneration, err)
	}

	user := &User{
		ID:           userID,
		Email:        input.Email,
		PasswordHash: string(bytes),
	}

	err = s.store.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, error) {
	// Normalize login input data
	email = strings.ToLower(strings.TrimSpace(email))

	v := validator.NewValidator()
	if ValidateUser(v, &input); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", ErrInvalidCredentials)
	}

	return user, nil
}
