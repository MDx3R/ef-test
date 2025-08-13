package dto

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	UserID      string     `json:"user_id" binding:"required,uuid"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
}

type SubscriptionQueryRequest struct {
	UserID      *string    `form:"user_id" binding:"omitempty,uuid"`
	ServiceName *string    `form:"service_name"`
	StartDate   *MonthYear `form:"start_date" example:"08-2025"`
	EndDate     *MonthYear `form:"end_date" example:"09-2025"`

	Page     int `form:"page,default=1,gte=1"`
	PageSize int `form:"page_size,default=20,gte=1"`
}

type TotalCostQueryRequest struct {
	UserID      string    `form:"user_id" binding:"required,uuid"`
	ServiceName string    `form:"service_name" binding:"required"`
	PeriodStart MonthYear `form:"period_start" binding:"required" example:"08-2025"`
	PeriodEnd   MonthYear `form:"period_end" binding:"required" example:"09-2025"`
}
