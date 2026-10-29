package auth

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/osugodbless/forum-backend.git/internal/auth/adapters/postgres"
	"github.com/osugodbless/forum-backend.git/internal/auth/adapters/restapi"
	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
	"github.com/osugodbless/forum-backend.git/internal/platform/jsonutil"
)

type Module struct {
	httpHandler *restapi.Handler
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Module {
	store := postgres.NewUserModel(db)

	authService := domain.NewAuthService(store)
	responder := jsonutil.NewResponder(logger)
	authHandler := restapi.NewHandler(authService, responder)

	return &Module{httpHandler: authHandler}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.httpHandler.RegisterAuthRoutes(mux)
}
