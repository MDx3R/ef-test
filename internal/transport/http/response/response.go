package response

import "github.com/google/uuid"

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

type IntResponse struct {
	Value int `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"validation error"`
	Fields map[string]string `json:"fields"`
}
