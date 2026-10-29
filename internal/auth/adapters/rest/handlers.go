package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
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
	service AuthService
	logger  *slog.Logger
}

func NewHandler(logger *slog.Logger, service AuthService) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var u SignupRequest

	err := h.readJSON(w, r, &u)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Normalize payload data
	u.FullName = strings.Join(strings.Fields(u.FullName), " ")
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	v := NewValidator()
	if ValidateUser(v, &u); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.service.SignUp(r.Context(), domain.SignupInput{
		FullName: u.FullName,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrUserAlreadyExists):
			h.requestConflictResponse(w, r, ErrUserAlreadyExists)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelopeJSON{"user": user}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
