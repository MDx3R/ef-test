package dto

import "github.com/google/uuid"

type IDResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type IntResponse struct {
	Value int `json:"value" example:"999"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"validation error"`
	Fields map[string]string `json:"fields"`
}

type SubscriptionResponse struct {
	ID          string     `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string     `json:"service_name" example:"Netflix"`
	Price       int        `json:"price" example:"999"`
	UserID      string     `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   MonthYear  `json:"start_date" example:"08-2025"`
	EndDate     *MonthYear `json:"end_date,omitempty" example:"09-2025"`
}
