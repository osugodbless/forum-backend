package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/osugodbless/forum-backend/internal/platform/validator"
	"golang.org/x/crypto/bcrypt"
)

const (
	// dummyHash is a bcrypt hash of a random password used for timing attack mitigation.
	dummyHash = "$2a$10$WkcuMITsGEtdm/ScvldHzekBfu9qXc0ANJZa3lcF7qvER0Nag2s1S"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")
)

type AuthService struct {
	store authStore
}

func NewAuthService(store authStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) SignUp(ctx context.Context, email, password string) (*User, error) {
	// Normalize signup data
	email = strings.ToLower(strings.TrimSpace(email))

	v := validator.NewValidator()
	if ValidateUser(v, email, password); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("sign up: generate user id: %w", err)
	}

	user := &User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.store.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth service: sign up: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, error) {
	// Normalize login data
	email = strings.ToLower(strings.TrimSpace(email))

	v := validator.NewValidator()
	if ValidateUser(v, email, password); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		// Mitigate timing attacks
		_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))

		return nil, fmt.Errorf("auth service: login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("auth service: login: %w", ErrInvalidCredentials)
	}

	return user, nil
}
