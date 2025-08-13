package gormmodel

import (
	"time"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/google/uuid"
)

type SubscriptionModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceName string
	Price       int
	UserID      uuid.UUID  `gorm:"type:uuid"`
	StartDate   time.Time  `gorm:"type:date"`
	EndDate     *time.Time `gorm:"type:date"`
}

func FromEntity(entity *entity.Subscription) SubscriptionModel {
	return SubscriptionModel{
		ID:          entity.ID(),
		ServiceName: entity.ServiceName(),
		Price:       entity.Price(),
		UserID:      entity.UserID(),
		StartDate:   entity.StartDate(),
		EndDate:     entity.EndDate(),
	}
}

func (m *SubscriptionModel) ToEntity() (*entity.Subscription, error) {
	sub, err := entity.NewSubscriptionWithID(
		m.ID,
		m.ServiceName,
		m.UserID,
		m.Price,
		m.StartDate,
		m.EndDate,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}
