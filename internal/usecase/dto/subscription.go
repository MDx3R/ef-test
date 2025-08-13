package dto

import (
	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/MDx3R/ef-test/internal/usecase/model"
	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID          uuid.UUID        `json:"id"`
	ServiceName string           `json:"service_name"`
	Price       int              `json:"price"`
	UserID      uuid.UUID        `json:"user_id"`
	StartDate   model.MonthYear  `json:"start_date"`
	EndDate     *model.MonthYear `json:"end_date,omitempty"`
}

type CreateSubscriptionRequests struct {
	ServiceName string           `json:"service_name" binding:"required"`
	Price       int              `json:"price" binding:"required"`
	UserID      uuid.UUID        `json:"user_id" binding:"required"`
	StartDate   model.MonthYear  `json:"start_date" binding:"required"`
	EndDate     *model.MonthYear `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequests struct {
	ServiceName string           `json:"service_name" binding:"required"`
	Price       int              `json:"price" binding:"required"`
	StartDate   model.MonthYear  `json:"start_date" binding:"required"`
	EndDate     *model.MonthYear `json:"end_date,omitempty"`
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID       `form:"user_id"`
	ServiceName *string          `form:"service_name"`
	StartDate   *model.MonthYear `form:"start_date"`
	EndDate     *model.MonthYear `form:"end_date"`

	Page     int `form:"page,default=1,gte=1"`
	PageSize int `form:"page_size,default=20,gte=1"`
}

type TotalCostFilter struct {
	UserID      uuid.UUID       `form:"user_id" binding:"required"`
	ServiceName string          `form:"service_name" binding:"required"`
	PeriodStart model.MonthYear `form:"period_start" binding:"required"`
	PeriodEnd   model.MonthYear `form:"period_end" binding:"required"`
}

func FromSubscription(sub *entity.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID(),
		ServiceName: sub.ServiceName(),
		Price:       sub.Price(),
		UserID:      sub.UserID(),
		StartDate:   model.NewMonthYear(sub.StartDate()),
		EndDate:     model.NewMonthYearFromPtr(sub.EndDate()),
	}
}
