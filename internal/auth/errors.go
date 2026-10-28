package auth

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")
)

func (h *Handler) logError(r *http.Request, err error) {
	h.logger.Error("http request failed", "error", err.Error(), "method", r.Method, "url", r.URL.RequestURI())
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := h.writeJSON(w, status, envelopeJSON{"error": message}, nil)
	if err != nil {
		h.logError(r, err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	h.errorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusNotFound, err.Error())
}

// func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
// 	h.errorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not allowed for the requested resource", r.Method))
// }

func (h *Handler) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (h *Handler) requestConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusConflict, err.Error())
}
