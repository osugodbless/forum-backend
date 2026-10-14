package main

import (
	"encoding/json"
	"net/http"

	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func routes(authSV *auth.Service) *http.ServeMux {
	mux := http.NewServeMux()
	authSV.RegisterAuthRoutes(mux)

	mux.HandleFunc("GET /v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(json.RawMessage(`{"status": "available", "environment": "development", "version": "1.0.0"}`))
	})
	return mux
}
