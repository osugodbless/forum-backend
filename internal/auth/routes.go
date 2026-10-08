package auth

import "net/http"

func AuthRoutes(mux *http.ServeMux, authSV *Service) {
	mux.HandleFunc("POST /register", authSV.handleSignUp)
	mux.HandleFunc("POST /login", authSV.handleLogin)
}
