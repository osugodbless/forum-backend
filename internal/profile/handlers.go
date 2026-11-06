package profile

import (
	"net/http"

	"github.com/osugodbless/forum-backend/internal/platform/jsonutil"
)

type Handler struct {
	service profileServiceManager
	*jsonutil.Responder
}

func NewHandler(service profileServiceManager, responder *jsonutil.Responder) *Handler {
	return &Handler{
		service:   service,
		Responder: responder,
	}
}

func (h *Handler) handleCreateProfile(w http.ResponseWriter, r *http.Request) {
	h.WriteJSON(w, http.StatusOK, "It worked", nil)
}
