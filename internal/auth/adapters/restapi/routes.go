package restapi

import "net/http"

func (h *Handler) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/register", h.handleSignUp)
	mux.HandleFunc("POST /v1/login", h.handleLogin)
}
