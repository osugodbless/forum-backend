package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) handleSignUp(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(`{"name": "Godbless", "email": "osugodbless234@gmail.com"}`)
}
