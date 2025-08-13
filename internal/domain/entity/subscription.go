package entity

import (
	"time"

	"github.com/MDx3R/ef-test/internal/domain"
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

func (s *Subscription) SetStartDate(startDate time.Time) error {
	if err := validateTime(startDate, s.endDate); err != nil {
		return err
	}

	s.startDate = startDate
	return nil
}

func (s *Subscription) SetEndDate(endDate *time.Time) error {
	if err := validateTime(s.startDate, endDate); err != nil {
		return err
	}

	s.endDate = endDate
	return nil
}

func (s *Subscription) SetStartEndDate(startDate time.Time, endDate *time.Time) error {
	if err := validateTime(startDate, endDate); err != nil {
		return err
	}

	s.startDate = startDate
	s.endDate = endDate
	return nil
}

func NewSubscription(
	serviceName string,
	userID uuid.UUID,
	price int,
	startDate time.Time,
	endDate *time.Time,
) (*Subscription, error) {
	if err := validateTime(startDate, endDate); err != nil {
		return nil, err
	}

	return &Subscription{
		id:          uuid.New(),
		serviceName: serviceName,
		price:       price,
		userID:      userID,
		startDate:   startDate,
		endDate:     endDate,
	}, nil
}

func NewSubscriptionWithID(
	id uuid.UUID,
	serviceName string,
	userID uuid.UUID,
	price int,
	startDate time.Time,
	endDate *time.Time,
) (*Subscription, error) {
	if err := validateTime(startDate, endDate); err != nil {
		return nil, err
	}

	return &Subscription{
		id:          id,
		serviceName: serviceName,
		price:       price,
		userID:      userID,
		startDate:   startDate,
		endDate:     endDate,
	}, nil
}

func validateTime(startDate time.Time, endDate *time.Time) error {
	if endDate != nil && startDate.After(*endDate) {
		return domain.ErrInvalidPeriod
	}
	return nil
}
