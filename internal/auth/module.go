package auth

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/osugodbless/forum-backend.git/internal/auth/adapters/postgres"
	"github.com/osugodbless/forum-backend.git/internal/auth/adapters/rest"
	"github.com/osugodbless/forum-backend.git/internal/auth/domain"
)

type Module struct {
	httpHandler *rest.Handler
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Module {
	store := postgres.NewUserModel(db)

	authService := domain.NewAuthService(store)
	authHandler := rest.NewHandler(logger, authService)

	return &Module{httpHandler: authHandler}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.httpHandler.RegisterAuthRoutes(mux)
}
