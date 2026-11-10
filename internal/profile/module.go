package profile

import (
	"log/slog"
	"net/http"

	"github.com/osugodbless/forum-backend/internal/platform/database"
	"github.com/osugodbless/forum-backend/internal/platform/jsonutil"
)

type Module struct {
	Service     *ProfileService
	httpHandler *Handler
}

type Deps struct {
	Logger *slog.Logger
	DB     database.DBTX
}

func New(d Deps) *Module {
	store := NewProfileModel(d.DB)

	profileService := NewProfileService(store)
	responder := jsonutil.NewResponder(d.Logger)
	profileHandler := NewHandler(profileService, responder)

	return &Module{
		Service:     profileService,
		httpHandler: profileHandler,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.httpHandler.RegisterProfileRoutes(mux)
}
