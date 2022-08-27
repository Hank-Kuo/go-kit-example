package response

type ErrorWrapper struct {
	Error string `json:"error"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
