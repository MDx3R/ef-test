package dto

import (
	"time"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CreateSubscriptionRequests struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequests struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID `form:"user_id"`
	ServiceName *string    `form:"service_name"`
	StartDate   *time.Time `form:"start_date"`
	EndDate     *time.Time `form:"end_date"`

	Page     int `form:"page,default=1,gte=1"`
	PageSize int `form:"page_size,default=20,gte=1"`
}

type TotalCostFilter struct {
	UserID      uuid.UUID `form:"user_id" binding:"required"`
	ServiceName string    `form:"service_name" binding:"required"`
	PeriodStart time.Time `form:"period_start" binding:"required"`
	PeriodEnd   time.Time `form:"period_end" binding:"required"`
}

func FromSubscription(sub *entity.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID(),
		ServiceName: sub.ServiceName(),
		Price:       sub.Price(),
		UserID:      sub.UserID(),
		StartDate:   sub.StartDate(),
		EndDate:     sub.EndDate(),
	}
}
