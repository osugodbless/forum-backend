package auth

import (
	"encoding/json"
	"net/http"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) handleSignUp(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(`{"name": "Godbless", "email": "osugodbless234@gmail.com"}`)
}
