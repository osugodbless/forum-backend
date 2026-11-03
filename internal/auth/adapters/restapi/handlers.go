package restapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
	"github.com/osugodbless/forum-backend.git/internal/platform/jsonutil"
	"github.com/osugodbless/forum-backend.git/internal/platform/validator"
)

type Handler struct {
	service authService
	*jsonutil.Responder
}

func NewHandler(service authService, responder *jsonutil.Responder) *Handler {
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

	user, err := h.service.SignUp(r.Context(), domain.SignupInput{
		FullName: u.FullName,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		var validationErr *validator.Validator
		switch {
		case errors.As(err, &validationErr):
			h.FailedValidationResponse(w, r, validationErr.Errors)
		case errors.Is(err, domain.ErrUserAlreadyExists):
			h.RequestConflictResponse(w, r, domain.ErrUserAlreadyExists)
		case errors.Is(err, domain.ErrIDGeneration):
			h.ServerErrorResponse(w, r, domain.ErrIDGeneration)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%s", user.ID.String()))

	err = h.WriteJSON(w, http.StatusCreated, toUserResponse(user), headers)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}

}

func toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
