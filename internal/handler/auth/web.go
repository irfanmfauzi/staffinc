package authHandler

import (
	"encoding/json"
	"net/http"
	"staffinc/internal/model/request"
	"staffinc/internal/model/response"
	"staffinc/internal/view/login"
	"staffinc/internal/view/register"
)

func (a *authHandler) GetRegisterWebHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-html")

	code := r.PathValue("code")

	register.Register(code).Render(r.Context(), w)
}

func (a *authHandler) PostRegisterWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	code := r.FormValue("code")

	errX := a.authService.Register(r.Context(), request.RegisterRequest{Email: email, Password: password, Role: "contributor"}, code)
	if errX.IsNotEmpty() {
		w.WriteHeader(errX.GetHttpCode())
		w.Write([]byte(errX.GetErrorCodeMessage().Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (a *authHandler) GetLoginWebHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text-html")

	login.Login().Render(r.Context(), w)
}

func (a *authHandler) LoginWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	req := request.LoginRequest{}
	req.Email = r.FormValue("email")
	req.Password = r.FormValue("password")

	tokenString, errX := a.authService.Login(r.Context(), req)
	if errX.IsNotEmpty() {
		w.Header().Set("Content-type", "application-json")
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: errX.GetErrorCodeMessage().Error()})
		w.WriteHeader(errX.GetHttpCode())
		w.Write(resp)
		return
	}

	cookie := http.Cookie{
		Name:  "staffinc_session",
		Value: tokenString,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
