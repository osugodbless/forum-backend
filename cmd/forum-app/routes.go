package main

import (
	"encoding/json"
	"net/http"

	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func routes(authSV *auth.Service) *http.ServeMux {
	mux := http.NewServeMux()
	authSV.RegisterAuthRoutes(mux)

	mux.HandleFunc("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		serverHealth := struct {
			Status     string            `json:"status"`
			SystemInfo map[string]string `json:"system_info"`
		}{
			Status: "available",
			SystemInfo: map[string]string{
				"environment": "development",
				"version":     "1.0.0",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(serverHealth); err != nil {
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}

	})

	return mux
}
