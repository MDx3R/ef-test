package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	id          uuid.UUID
	serviceName string
	price       int
	userID      uuid.UUID
	startDate   time.Time
	endDate     *time.Time
}

func (s *Subscription) ID() uuid.UUID {
	return s.id
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

func (s *Subscription) EndDate() *time.Time {
	return s.endDate
}

func (s *Subscription) SetServiceName(serviceName string) {
	s.serviceName = serviceName
}

func (s *Subscription) SetPrice(price int) {
	s.price = price
}

func (s *Subscription) SetStartDate(startDate time.Time) {
	s.startDate = startDate
}

func (s *Subscription) SetEndDate(endDate *time.Time) {
	s.endDate = endDate
}

func NewSubscription(
	serviceName string,
	userID uuid.UUID,
	price int,
	startDate time.Time,
	endDate *time.Time,
) *Subscription {
	return &Subscription{
		id:          uuid.New(),
		serviceName: serviceName,
		price:       price,
		userID:      userID,
		startDate:   startDate,
		endDate:     endDate,
	}
}
