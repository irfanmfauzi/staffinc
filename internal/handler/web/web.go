package webHandler

import (
	"net/http"
	generatorlinkService "staffinc/internal/service/generator_link"
	"staffinc/internal/view/dashboard"
	"staffinc/middleware"
)

type webHandler struct {
	generatorlinkService generatorlinkService.GeneratorLinkServiceProvider
}

func RegisterWebHandlerRoute(mux *http.ServeMux, service generatorlinkService.GeneratorLinkServiceProvider) {
	handler := &webHandler{
		generatorlinkService: service,
	}
	mux.Handle("GET /dashboard", middleware.VerifyToken(http.HandlerFunc(handler.Dashboard)))
}

func (h *webHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	items, err := h.generatorlinkService.GetLink(r.Context())
	if err != nil {
		return
	}
	dashboard.Dashboard(items).Render(r.Context(), w)
}
