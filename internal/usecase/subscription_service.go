package usecase

import (
	"errors"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetSubscription(id uuid.UUID) (dto.SubscriptionResponse, error)
	ListSubscriptions(filter dto.SubscriptionFilter) ([]dto.SubscriptionResponse, error)
	CreateSubscription(request dto.CreateSubscriptionRequests) (uuid.UUID, error)
	UpdateSubscription(id uuid.UUID, request dto.UpdateSubscriptionRequests) error
	DeleteSubscription(id uuid.UUID) error
	CalculateTotalCost(filter dto.TotalCostFilter) (int, error)
}

type subscriptionService struct {
	subRepo SubscriptionRepository
}

func NewSubscriptionService(subRepo SubscriptionRepository) SubscriptionService {
	return &subscriptionService{subRepo: subRepo}
}

func (s *subscriptionService) GetSubscription(id uuid.UUID) (dto.SubscriptionResponse, error) {
	sub, err := s.subRepo.Get(id)
	if err != nil {
		return dto.SubscriptionResponse{}, err
	}
	return dto.FromSubscription(sub), nil
}

func (s *subscriptionService) ListSubscriptions(filter dto.SubscriptionFilter) ([]dto.SubscriptionResponse, error) {
	subs, err := s.subRepo.List(filter)
	if err != nil {
		return []dto.SubscriptionResponse{}, err
	}

	result := make([]dto.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		result[i] = dto.FromSubscription(sub)
	}

	return result, nil
}

func (s *subscriptionService) CreateSubscription(request dto.CreateSubscriptionRequests) (uuid.UUID, error) {
	sub, err := entity.NewSubscription(
		request.ServiceName,
		request.UserID,
		request.Price,
		request.StartDate.ToTime(),
		request.EndDate.ToTimePtr(),
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.subRepo.Add(sub); err != nil {
		return uuid.Nil, err
	}
	return sub.ID(), nil
}

func (s *subscriptionService) UpdateSubscription(id uuid.UUID, request dto.UpdateSubscriptionRequests) error {
	sub, err := s.subRepo.Get(id)
	if err != nil {
		return err
	}

	sub.SetServiceName(request.ServiceName)
	sub.SetPrice(request.Price)
	if err := sub.SetStartEndDate(request.StartDate.ToTime(), request.EndDate.ToTimePtr()); err != nil {
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
