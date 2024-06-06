package response

type GenericResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	GenericResponse
	Data TokenResponse `json:"data"`
}
