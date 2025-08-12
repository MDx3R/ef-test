package usecase

import (
	"time"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	UserID      uuid.UUID
	ServiceName string
	PeriodStart time.Time
	PeriodEnd   *time.Time
	Page        int
	PageSize    int
}

type SubscriptionRepository interface {
	Get(id uuid.UUID) (*entity.Subscription, error)
	List(filter SubscriptionFilter) ([]*entity.Subscription, error)
	Add(sub *entity.Subscription) error
	Update(sub *entity.Subscription) error
	Delete(id uuid.UUID) error
	CalculateTotalCost(filter SubscriptionFilter) (int, error)
}
