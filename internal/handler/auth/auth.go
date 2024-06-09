package authHandler

import (
	"net/http"
	authService "staffinc/internal/service/auth"
)

type authHandler struct {
	authService authService.AuthServiceProvider
}

func RegisterAuthRoute(mux *http.ServeMux, authService authService.AuthServiceProvider) {
	handler := &authHandler{
		authService: authService,
	}

	registerAPIAuth(mux, handler)
	registerWebAuth(mux, handler)
}

func registerAPIAuth(mux *http.ServeMux, handler *authHandler) {
	mux.HandleFunc("POST /api/auth/register", handler.RegisterApiHandler)
	mux.HandleFunc("POST /api/auth/register/{code}", handler.RegisterApiHandler)

	mux.HandleFunc("POST /api/auth/login", handler.LoginApiHandler)
}

func registerWebAuth(mux *http.ServeMux, handler *authHandler) {
	mux.HandleFunc("GET /register/{code}", handler.GetRegisterWebHandler)
	mux.HandleFunc("GET /register", handler.GetRegisterWebHandler)
	mux.HandleFunc("POST /register", handler.PostRegisterWebHandler)

	mux.HandleFunc("GET /login", handler.GetLoginWebHandler)
	mux.HandleFunc("POST /login", handler.LoginWebHandler)
}
