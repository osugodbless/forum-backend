package profile

import "net/http"

func (h *Handler) RegisterProfileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/profiles", h.handleCreateProfile)
}
