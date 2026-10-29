package restapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
	"github.com/osugodbless/forum-backend.git/internal/platform/jsonutil"
	"github.com/osugodbless/forum-backend.git/internal/platform/validator"
)

type Handler struct {
	service AuthService
	*jsonutil.Responder
}

func NewHandler(service AuthService, responder *jsonutil.Responder) *Handler {
	return &Handler{
		service:   service,
		Responder: responder,
	}
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var u SignupRequest

	err := h.ReadJSON(w, r, &u)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	// Normalize payload data
	u.FullName = strings.Join(strings.Fields(u.FullName), " ")
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	v := validator.NewValidator()
	if ValidateUser(v, &u); !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.service.SignUp(r.Context(), domain.SignupInput{
		FullName: u.FullName,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			h.RequestConflictResponse(w, r, domain.ErrUserAlreadyExists)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = h.WriteJSON(w, http.StatusCreated, jsonutil.EnvelopeJSON{"user": user}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
