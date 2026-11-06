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

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")
)

type AuthService struct {
	store   authStore
	profile ProfileCreator
}

func NewAuthService(store authStore, profile ProfileCreator) *AuthService {
	return &AuthService{
		store:   store,
		profile: profile,
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
	if ValidateUser(v, input.FullName, input.Email, input.Password); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("user id generation failed")
	}

	user := &User{
		ID:           userID,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	tx, err := s.store.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting registration transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	s.store.WithTx(tx)

	err = s.store.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, error) {
	// Normalize login input data
	email = strings.ToLower(strings.TrimSpace(email))

	v := validator.NewValidator()
	if ValidateUser(v, "", email, password); !v.Valid() {
		return nil, validator.NewValidatorError(v.Errors)
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("auth service: %w", ErrInvalidCredentials)
	}

	return user, nil
}
