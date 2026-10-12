package auth

import "net/http"

func (s *Service) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", s.handleSignUp)
	mux.HandleFunc("POST /login", s.handleLogin)
}
