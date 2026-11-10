package auth

import (
	"log/slog"
	"net/http"

	"github.com/osugodbless/forum-backend/internal/auth/adapters/postgres"
	"github.com/osugodbless/forum-backend/internal/auth/adapters/restapi"
	"github.com/osugodbless/forum-backend/internal/auth/domain"
	"github.com/osugodbless/forum-backend/internal/platform/database"
	"github.com/osugodbless/forum-backend/internal/platform/jsonutil"
)

type Module struct {
	httpHandler *restapi.Handler
}

type Deps struct {
	Logger  *slog.Logger
	DB      database.DBTX
	Tx      domain.TxManager
	Profile domain.ProfileCreator
}

func New(d Deps) *Module {
	store := postgres.NewUserModel(d.DB)

	authService := domain.NewAuthService(store, d.Tx, d.Profile)
	responder := jsonutil.NewResponder(d.Logger)
	authHandler := restapi.NewHandler(authService, responder)

	return &Module{httpHandler: authHandler}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.httpHandler.RegisterAuthRoutes(mux)
}
