package auth

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("username or email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrResourceNotFound   = errors.New("the requested resource could not be found")
)

func (s *Service) logError(r *http.Request, err error) {
	s.logger.Error(err.Error(), "method", r.Method, "url", r.URL.RequestURI())
}

func (s *Service) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := s.writeJSON(w, status, envelopeJSON{"error": message}, nil)
	if err != nil {
		s.logError(r, err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (s *Service) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)
	s.errorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (s *Service) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (s *Service) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusNotFound, err.Error())
}

func (s *Service) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	s.errorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not allowed for the requested resource", r.Method))
}
