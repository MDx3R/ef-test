package usecase_test

import (
	"testing"
	"time"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	mock_usecase "github.com/MDx3R/ef-test/internal/usecase/mocks"
	"github.com/MDx3R/ef-test/internal/usecase/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSubscriptionService(t *testing.T) (*mock_usecase.MockSubscriptionRepository, usecase.SubscriptionService) {
	mockRepo := mock_usecase.NewMockSubscriptionRepository(t)
	service := usecase.NewSubscriptionService(mockRepo)
	return mockRepo, service
}

func makeTestSubscription(t *testing.T) *entity.Subscription {
	id := uuid.New()
	sub, _ := entity.NewSubscriptionWithID(
		id,
		"test_service",
		uuid.New(),
		100,
		time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		nil,
	)

	return sub
}

func TestSubscriptionService_GetSubscription(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	sub := makeTestSubscription(t)
	id := sub.ID()

	mockRepo.On("Get", id).Return(sub, nil)

	resp, err := service.GetSubscription(id)

	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
	assert.Equal(t, "test_service", resp.ServiceName)
	assert.Equal(t, 100, resp.Price)
	assert.Equal(t, model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)), resp.StartDate)
	assert.Nil(t, resp.EndDate)
	assert.IsType(t, uuid.UUID{}, resp.UserID)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_GetSubscription_NotFound(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	sub := makeTestSubscription(t)
	id := sub.ID()

	mockRepo.On("Get", id).Return(nil, usecase.ErrNotFound)

	_, err := service.GetSubscription(id)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_ListSubscriptions(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	subs := []*entity.Subscription{
		makeTestSubscription(t),
		makeTestSubscription(t),
	}

	filter := dto.SubscriptionFilter{}

	mockRepo.On("List", filter).Return(subs, nil)

	resp, err := service.ListSubscriptions(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, subs[0].ID(), resp[0].ID)
	assert.Equal(t, subs[1].ID(), resp[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_ListSubscriptions_Error(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	filter := dto.SubscriptionFilter{}

	mockRepo.On("List", filter).Return(nil, usecase.ErrRepository)

	resp, err := service.ListSubscriptions(filter)

	assert.Error(t, err)
	assert.Empty(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscriptions(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	endDate := model.NewMonthYear(time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC))
	req := dto.CreateSubscriptionRequests{
		ServiceName: "service_test",
		Price:       100,
		StartDate:   model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:     &endDate,
	}

	mockRepo.On("Add", mock.AnythingOfType("*entity.Subscription")).Return(nil)

	id, err := service.CreateSubscription(req)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscriptions_Error(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	req := dto.CreateSubscriptionRequests{}

	mockRepo.On("Add", mock.AnythingOfType("*entity.Subscription")).Return(usecase.ErrRepository)

	id, err := service.CreateSubscription(req)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpdateSubscriptions(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	sub := makeTestSubscription(t)
	id := sub.ID()
	endDate := model.NewMonthYear(time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC))

	req := dto.UpdateSubscriptionRequests{
		ServiceName: "updated_name",
		Price:       150,
		StartDate:   model.NewMonthYear(time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:     &endDate,
	}

	mockRepo.On("Get", id).Return(sub, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *entity.Subscription) bool {
		return (u.ID() == sub.ID() &&
			u.ServiceName() == req.ServiceName &&
			u.UserID() == sub.UserID() &&
			time.Time.Equal(u.StartDate(), req.StartDate.ToTime()) &&
			time.Time.Equal(*u.EndDate(), req.EndDate.ToTime()))
	})).Return(nil)

	err := service.UpdateSubscription(id, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpdateSubscriptions_GetError(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	sub := makeTestSubscription(t)
	id := sub.ID()

	req := dto.UpdateSubscriptionRequests{}

	mockRepo.On("Get", id).Return(nil, usecase.ErrNotFound)

	err := service.UpdateSubscription(id, req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpdateSubscriptions_UpdateError(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	sub := makeTestSubscription(t)
	id := sub.ID()

	req := dto.UpdateSubscriptionRequests{}

	mockRepo.On("Get", id).Return(sub, nil)
	mockRepo.On("Update", mock.Anything).Return(usecase.ErrRepository)

	err := service.UpdateSubscription(id, req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeleteSubscription(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	id := uuid.New()

	mockRepo.On("Delete", id).Return(nil)

	err := service.DeleteSubscription(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeleteSubscription_Error(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	id := uuid.New()

	mockRepo.On("Delete", id).Return(usecase.ErrRepository)

	err := service.DeleteSubscription(id)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeleteSubscription_NotFound_Skips(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	id := uuid.New()

	mockRepo.On("Delete", id).Return(usecase.ErrNotFound)

	err := service.DeleteSubscription(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CalculateTotalCost(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	filter := dto.TotalCostFilter{
		UserID:      uuid.New(),
		ServiceName: "service_test",
		PeriodStart: model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)),
		PeriodEnd:   model.NewMonthYear(time.Date(2030, 9, 1, 0, 0, 0, 0, time.UTC)),
	}

	mockRepo.On("CalculateTotalCost", filter).Return(150, nil)

	result, err := service.CalculateTotalCost(filter)

	assert.NoError(t, err)
	assert.Equal(t, 150, result)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CalculateTotalCost_Error(t *testing.T) {
	mockRepo, service := setupSubscriptionService(t)

	filter := dto.TotalCostFilter{}

	mockRepo.On("CalculateTotalCost", filter).Return(0, usecase.ErrRepository)

	result, err := service.CalculateTotalCost(filter)

	assert.Error(t, err)
	assert.Equal(t, 0, result)
	mockRepo.AssertExpectations(t)
}
