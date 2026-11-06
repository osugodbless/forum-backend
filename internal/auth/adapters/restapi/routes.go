package restapi

import "net/http"

func (h *Handler) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/users", h.handleSignUp)
	mux.HandleFunc("POST /v1/login", h.handleLogin)
}
