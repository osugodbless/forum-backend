package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	cfg, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	authService := auth.NewService(logger)

	mux := routes(authService)

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.server.host, cfg.server.port),
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", slog.String("addr", server.Addr))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
