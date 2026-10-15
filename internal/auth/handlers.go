package auth

import (
	"encoding/json"
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
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) handleSignUp(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := json.RawMessage(`{"name": "Godbless", "email": "osugodbless234@gmail.com"}`)
	err := s.writeJSON(w, http.StatusOK, envelopeJSON{"user": data}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
