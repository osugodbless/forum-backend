package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/osugodbless/forum-backend/internal/platform/validator"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func ValidateUser(v *validator.Validator, email, password string) {
	// Email validation
	v.Check(email != "" && len(email) <= 254, "email", "Email field empty or too long")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "Email address is not valid")

	// Password validation
	v.Check(password != "", "password", "Password field is empty")
	v.Check(len(password) >= 8 && len(password) <= 64, "password", "Password must be between 8 and 64 characters long")
}
