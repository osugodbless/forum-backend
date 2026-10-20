package auth

import (
	"log/slog"
	"net/http"
)

// envelopeJSON envelopes response JSON data in a parent string
//
// Example a type User struct{"id": 123, "email": "example@gmail.com",} is enveloped in the "user" string as such;
//
//	{
//			"user": {
//				"id": 123,
//				"email": "example@gmail.com"
//				...
//			}
//	}
//
// Usage: envelopeJSON{"user": data}, where data in this instance is the user model
type envelopeJSON map[string]any

type Service struct {
	userModel *UserModel
	logger    *slog.Logger
}

func NewService(logger *slog.Logger, usermodel *UserModel) *Service {
	return &Service{
		logger:    logger,
		userModel: usermodel,
	}
}

func (s *Service) handleSignUp(w http.ResponseWriter, r *http.Request) {

}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	var u *User
	if err := s.readJSON(w, r, &u); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	err := s.writeJSON(w, http.StatusOK, envelopeJSON{"user": u}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
