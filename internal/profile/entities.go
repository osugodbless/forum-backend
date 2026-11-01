package profile

type UserProfile struct {
	UserID      string `json:"user_id" db:"user_id"`
	Username    string `json:"username" db:"username"`
	FullName    string `json:"full_name" db:"full_name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Country     string `json:"country" db:"country"`
	Bio         string `json:"bio" db:"bio"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}
