package authHandler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"staffinc/internal/model/request"
	"staffinc/internal/model/response"
)

// RegisterApiHandler is handler to create an event
// RegisterApiHandler godoc
// @Summary     RegisterApiHandler
// @Accept      json
// @Description Register, PIC : Irfan Fauzi
// @Produce     json
// @Tags        Auth
// @Param       RegisterRequest body     request.RegisterRequest true "RegisterRequest"
// @Success     201             {object} response.GenericResponse{}
// @Failure     500             {object} response.GenericResponse{}
// @ID          v1-RegisterApiHandler
// @Router      /auth/register   [post]
func (a *authHandler) RegisterApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	req := request.RegisterRequest{}
	req.Role = "generator"

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	ctx := r.Context()
	code := r.PathValue("code")
	httpCode, err := a.authService.Register(ctx, req, code)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(httpCode)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Register"})
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// LoginApiHandler is handler to create an event
// LoginApiHandler godoc
// @Summary     LoginApiHandler
// @Accept      json
// @Description Register, PIC : Irfan Fauzi
// @Produce     json
// @Tags        Auth
// @Param       RegisterRequest body     request.LoginRequest true "LoginRequest"
// @Success     201             {object} response.LoginResponse{}
// @Failure     500             {object} response.GenericResponse{}
// @ID          v1-LoginApiHandler
// @Router      /auth/login   [post]
func (a *authHandler) LoginApiHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application-json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Fail when reading body", "Error", err)
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}
	defer r.Body.Close()
	req := request.LoginRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	ctx := r.Context()

	tokenString, code, err := a.authService.Login(ctx, req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(code)
		w.Write(resp)
		return
	}

	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:  "staffinc_session",
		Value: tokenString,
	}
	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	resp, _ := json.Marshal(
		response.LoginResponse{
			GenericResponse: response.GenericResponse{Success: true, Message: "Login Success"},
			Data:            response.TokenResponse{Token: tokenString},
		},
	)

	w.Write(resp)
}
