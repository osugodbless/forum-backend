package profile

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Username    string    `json:"username" db:"username"`
	FullName    string    `json:"full_name" db:"full_name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Country     string    `json:"country" db:"country"`
	Bio         string    `json:"bio" db:"bio"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// // Full name validation
// 	v.Check(fullName != "", "full_name", "Full name field is empty")
// 	v.Check(utf8.RuneCountInString(fullName) <= 100, "full_name", "Full name must be less than or equal to 100 characters long")
