// Example User Model
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUser(v *Validator, user *User) {

}