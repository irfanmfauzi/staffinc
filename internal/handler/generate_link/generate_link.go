package generatelinkHandler

import (
	"encoding/json"
	"net/http"
	"staffinc/internal/model/response"
	generatorlinkService "staffinc/internal/service/generator_link"
	"staffinc/middleware"
)

type generateLinkHandler struct {
	generateLinkService generatorlinkService.GeneratorLinkServiceProvider
}

func RegisterGenerateLink(mux *http.ServeMux, service generatorlinkService.GeneratorLinkServiceProvider) {
	handler := &generateLinkHandler{
		generateLinkService: service,
	}

	mux.Handle("POST /api/generate-link", middleware.VerifyToken(http.HandlerFunc(handler.PostGenerateLinkHandler)))

}

// PostGenerateLinkHandler is handler to create an event
// PostGenerateLinkHandler godoc
// @Summary     PostGenerateLinkHandler
// @Accept      json
// @Description Register, PIC : Irfan Fauzi
// @Produce     json
// @Tags        Generate Link
// @Success     201 {object} response.GenericResponse{}
// @Failure     500 {object} response.GenericResponse{}
// @ID          v1-PostGenerateLinkHandler
// @Router      /generic-response   [post]
func (h *generateLinkHandler) PostGenerateLinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	ctx := r.Context()
	user := ctx.Value("user").(map[string]interface{})
	userId := int64(user["id"].(float64))
	role := user["role"].(string)

	errX := h.generateLinkService.GenerateLink(ctx, userId, role)
	if errX.IsNotEmpty() {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: errX.GetErrorCodeMessage().Error()})
		w.WriteHeader(errX.GetHttpCode())
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success generate link"})
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	return
}
