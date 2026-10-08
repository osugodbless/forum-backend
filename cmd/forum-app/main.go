package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	authService := auth.NewService()

	mux := routes(authService)

	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
