package handler

import (
	"net/http"
	authHandler "staffinc/internal/handler/auth"
	webHandler "staffinc/internal/handler/web"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	authHandler.RegisterAuthRoute(mux, s.service.UserService)
	webHandler.RegisterWebHandlerRoute(mux)
	return mux
}
