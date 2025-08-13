package dto

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required" example:"Netflix"`
	Price       int        `json:"price" binding:"required" example:"999"`
	UserID      string     `json:"user_id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   MonthYear  `json:"start_date" binding:"required" example:"08-2025"`
	EndDate     *MonthYear `json:"end_date,omitempty" example:"09-2025"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required" example:"Netflix"`
	Price       int        `json:"price" binding:"required" example:"999"`
	StartDate   MonthYear  `json:"start_date" binding:"required" example:"08-2025"`
	EndDate     *MonthYear `json:"end_date,omitempty" example:"09-2025"`
}

type SubscriptionQueryRequest struct {
	UserID      *string    `form:"user_id" binding:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName *string    `form:"service_name" example:"Netflix"`
	StartDate   *MonthYear `form:"start_date" example:"08-2025"`
	EndDate     *MonthYear `form:"end_date" example:"09-2025"`

	Page     int `form:"page,default=1,gte=1" example:"1"`
	PageSize int `form:"page_size,default=20,gte=1" example:"20"`
}

type TotalCostQueryRequest struct {
	UserID      string    `form:"user_id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string    `form:"service_name" binding:"required" example:"Netflix"`
	PeriodStart MonthYear `form:"period_start" binding:"required" example:"08-2025"`
	PeriodEnd   MonthYear `form:"period_end" binding:"required" example:"09-2025"`
}
