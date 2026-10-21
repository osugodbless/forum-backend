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

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(logger *slog.Logger, service *Service) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var u SignupDTO

	err := h.readJSON(w, r, &u)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
