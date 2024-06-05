package response

type GenericResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}
