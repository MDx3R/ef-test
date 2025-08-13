package dto

import (
	"time"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/google/uuid"
)

type SubscriptionDTO struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type CreateSubscriptionCommand struct {
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type UpdateSubscriptionCommand struct {
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time

	Page     int
	PageSize int
}

type TotalCostFilter struct {
	UserID      uuid.UUID
	ServiceName string
	PeriodStart time.Time
	PeriodEnd   time.Time
}

func FromSubscription(sub *entity.Subscription) SubscriptionDTO {
	return SubscriptionDTO{
		ID:          sub.ID(),
		ServiceName: sub.ServiceName(),
		Price:       sub.Price(),
		UserID:      sub.UserID(),
		StartDate:   sub.StartDate(),
		EndDate:     sub.EndDate(),
	}
}
