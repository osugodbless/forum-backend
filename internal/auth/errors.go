package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("username or email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
