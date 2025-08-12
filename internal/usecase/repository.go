package usecase

import (
	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Get(id uuid.UUID) (*entity.Subscription, error)
	List(filter dto.SubscriptionFilter) ([]*entity.Subscription, error)
	Add(sub *entity.Subscription) error
	Update(sub *entity.Subscription) error
	Delete(id uuid.UUID) error
	CalculateTotalCost(filter dto.TotalCostFilter) (int, error)
}
