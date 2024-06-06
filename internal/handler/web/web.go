package webHandler

import (
	"net/http"
	"staffinc/internal/view/layout/admin"
	"staffinc/middleware"
)

type webHandler struct{}

func RegisterWebHandlerRoute(mux *http.ServeMux) {
	handler := &webHandler{}
	mux.Handle("GET /dashboard", middleware.VerifyToken(http.HandlerFunc(handler.Dashboard)))
}

func (h *webHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	admin.Base().Render(r.Context(), w)
}
