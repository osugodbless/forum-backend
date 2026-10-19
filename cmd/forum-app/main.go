package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	cfg, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	pool, err := OpenDB(cfg.DB.URL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("Database connected sucessfully")

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

func OpenDB(dsn string) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to create database pool: %s\n", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Unable to ping database: %s\n", err)
	}

	return pool, nil
}
