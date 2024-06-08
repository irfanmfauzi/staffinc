package handler

import (
	"net/http"
	_ "staffinc/docs"
	authHandler "staffinc/internal/handler/auth"
	generatelinkHandler "staffinc/internal/handler/generate_link"
	webHandler "staffinc/internal/handler/web"

	"github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.Handler())
	authHandler.RegisterAuthRoute(mux, s.service.UserService)
	webHandler.RegisterWebHandlerRoute(mux, s.service.GenerateLinkService)
	generatelinkHandler.RegisterGenerateLink(mux, s.service.GenerateLinkService)

	return mux
}
