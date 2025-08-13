package usecase

import (
	"errors"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetSubscription(id uuid.UUID) (dto.SubscriptionDTO, error)
	ListSubscriptions(filter dto.SubscriptionFilter) ([]dto.SubscriptionDTO, error)
	CreateSubscription(request dto.CreateSubscriptionCommand) (uuid.UUID, error)
	UpdateSubscription(id uuid.UUID, request dto.UpdateSubscriptionCommand) error
	DeleteSubscription(id uuid.UUID) error
	CalculateTotalCost(filter dto.TotalCostFilter) (int, error)
}

type subscriptionService struct {
	subRepo SubscriptionRepository
}

func NewSubscriptionService(subRepo SubscriptionRepository) SubscriptionService {
	return &subscriptionService{subRepo: subRepo}
}

func (s *subscriptionService) GetSubscription(id uuid.UUID) (dto.SubscriptionDTO, error) {
	sub, err := s.subRepo.Get(id)
	if err != nil {
		return dto.SubscriptionDTO{}, err
	}
	return dto.FromSubscription(sub), nil
}

func (s *subscriptionService) ListSubscriptions(filter dto.SubscriptionFilter) ([]dto.SubscriptionDTO, error) {
	subs, err := s.subRepo.List(filter)
	if err != nil {
		return []dto.SubscriptionDTO{}, err
	}

	result := make([]dto.SubscriptionDTO, len(subs))
	for i, sub := range subs {
		result[i] = dto.FromSubscription(sub)
	}

	return result, nil
}

func (s *subscriptionService) CreateSubscription(request dto.CreateSubscriptionCommand) (uuid.UUID, error) {
	sub, err := entity.NewSubscription(
		request.ServiceName,
		request.UserID,
		request.Price,
		request.StartDate,
		request.EndDate,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.subRepo.Add(sub); err != nil {
		return uuid.Nil, err
	}
	return sub.ID(), nil
}

func (s *subscriptionService) UpdateSubscription(id uuid.UUID, request dto.UpdateSubscriptionCommand) error {
	sub, err := s.subRepo.Get(id)
	if err != nil {
		return err
	}

	sub.SetServiceName(request.ServiceName)
	sub.SetPrice(request.Price)
	if err := sub.SetStartEndDate(request.StartDate, request.EndDate); err != nil {
		return err
	}

	if err := s.subRepo.Update(sub); err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) DeleteSubscription(id uuid.UUID) error {
	err := s.subRepo.Delete(id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	return nil
}

func (s *subscriptionService) CalculateTotalCost(filter dto.TotalCostFilter) (int, error) {
	return s.subRepo.CalculateTotalCost(filter)
}
