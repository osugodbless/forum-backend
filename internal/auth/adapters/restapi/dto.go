package restapi

import (
	"unicode/utf8"

	"github.com/osugodbless/forum-backend.git/internal/platform/validator"
)

// DTOs (Data Transfer Object)
type SignupRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUser(v *validator.Validator, u *SignupRequest) {
	// Email validation
	v.Check(u.Email != "" && len(u.Email) <= 254, "email", "Email field empty or too long")
	v.Check(validator.Matches(u.Email, validator.EmailRX), "email", "Email address is not valid")

	// Password validation
	v.Check(u.Password != "", "password", "Password field is empty")
	v.Check(len(u.Password) >= 8 && len(u.Password) <= 64, "password", "Password must be between 8 and 64 characters long")

	// Full name validation
	v.Check(u.FullName != "", "full_name", "Full name field is empty")
	v.Check(utf8.RuneCountInString(u.FullName) <= 100, "full_name", "Full name must be less than or equal to 100 characters long")
}
