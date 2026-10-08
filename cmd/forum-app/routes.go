package main

import (
	"net/http"

	"github.com/osugodbless/forum-backend.git/internal/auth"
)

func routes(authSV *auth.Service) *http.ServeMux {
	mux := http.NewServeMux()
	auth.AuthRoutes(mux, authSV)
	return mux
}
