package profile

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/osugodbless/forum-backend.git/internal/platform/jsonutil"
)

type Module struct {
	httpHandler *Handler
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Module {
	store := NewProfileModel(db)

	profileService := NewProfileService(store)
	responder := jsonutil.NewResponder(logger)
	profileHandler := NewHandler(profileService, responder)

	return &Module{httpHandler: profileHandler}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.httpHandler.RegisterProfileRoutes(mux)
}
