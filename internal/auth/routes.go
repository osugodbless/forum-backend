package auth

import "net/http"

func (s *Service) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/register", s.handleSignUp)
	mux.HandleFunc("POST /v1/login", s.handleLogin)
}
