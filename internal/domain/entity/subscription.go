package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	serviceName string
	price       int
	userID      uuid.UUID
	startDate   time.Time
	endDate     *time.Time
}

func (s *Subscription) ServiceName() string {
	return s.serviceName
}

func (s *Subscription) Price() int {
	return s.price
}

func (s *Subscription) UserID() uuid.UUID {
	return s.userID
}

func (s *Subscription) StartDate() time.Time {
	return s.startDate
}

func (s *Subscription) EndDate() time.Time {
	return *s.endDate
}

func NewSubscription(
	serviceName string,
	userID uuid.UUID,
	price int,
	startDate time.Time,
	endDate *time.Time,
) *Subscription {
	return &Subscription{
		serviceName: serviceName,
		price:       price,
		userID:      userID,
		startDate:   startDate,
		endDate:     endDate,
	}
}
