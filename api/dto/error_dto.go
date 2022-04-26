package dto

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
