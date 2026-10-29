package rest

import (
	"regexp"
	"unicode/utf8"
)

// EmailRX is a regular expression for sanity checking the format of email addresses. The regular expression pattern is
// taken from https://html.spec.whatwg.org/#valid-e-mail-address.
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")
)

// Validator is a custom type which holds a map of validation errors.
// The keys in the map are the names of the fields which failed validation, and the values are the corresponding error messages.
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance with an empty errors map.
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a specific value is in a list of permitted values and returns the corresponding boolean value.
// func PermittedValue[T comparable](value T, permittedValues ...T) bool {
// 	return slices.Contains(permittedValues, value)
// }

// Matches checks if a string value matches a specific regexp pattern, then returns the corresponding boolean value.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all values in a slice are unique and returns the corresponding boolean value.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}

func ValidateUser(v *Validator, u *SignupRequest) {
	// Email validation
	v.Check(u.Email != "" && len(u.Email) <= 254, "email", "Email field empty or too long")
	v.Check(Matches(u.Email, EmailRX), "email", "Email address is not valid")

	// Password validation
	v.Check(u.Password != "", "password", "Password field is empty")
	v.Check(len(u.Password) >= 8 && len(u.Password) <= 64, "password", "Password must be between 8 and 64 characters long")

	// Full name validation
	v.Check(u.FullName != "", "full_name", "Full name field is empty")
	v.Check(utf8.RuneCountInString(u.FullName) <= 100, "full_name", "Full name must be less than or equal to 100 characters long")
}
