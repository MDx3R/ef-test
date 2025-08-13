package dto

import (
	"time"

	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/google/uuid"
)

func ToCreateSubscriptionCommand(r CreateSubscriptionRequest) (*dto.CreateSubscriptionCommand, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}

	startDate, err := r.StartDate.Parse()
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	endDate, err = r.EndDate.Parse()
	if err != nil {
		return nil, err
	}

	return &dto.CreateSubscriptionCommand{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   *startDate,
		EndDate:     endDate,
	}, nil
}

func ToUpdateSubscriptionCommand(r UpdateSubscriptionRequest) (*dto.UpdateSubscriptionCommand, error) {
	startDate, err := r.StartDate.Parse()
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	endDate, err = r.EndDate.Parse()
	if err != nil {
		return nil, err
	}

	return &dto.UpdateSubscriptionCommand{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		StartDate:   *startDate,
		EndDate:     endDate,
	}, nil
}

func ToSubscriptionFilter(r SubscriptionQueryRequest) (*dto.SubscriptionFilter, error) {
	var userID *uuid.UUID
	if r.UserID != nil {
		uid, err := uuid.Parse(*r.UserID)
		if err != nil {
			return nil, err
		}
		userID = &uid
	}

	var startDate, endDate *time.Time
	var err error
	startDate, err = r.StartDate.Parse()
	if err != nil {
		return nil, err
	}
	endDate, err = r.EndDate.Parse()
	if err != nil {
		return nil, err
	}

	return &dto.SubscriptionFilter{
		UserID:      userID,
		ServiceName: r.ServiceName,
		StartDate:   startDate,
		EndDate:     endDate,
		Page:        r.Page,
		PageSize:    r.PageSize,
	}, nil
}

func ToTotalCostFilter(r TotalCostQueryRequest) (*dto.TotalCostFilter, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}

	periodStart, err := r.PeriodStart.Parse()
	if err != nil {
		return nil, err
	}

	periodEnd, err := r.PeriodEnd.Parse()
	if err != nil {
		return nil, err
	}

	return &dto.TotalCostFilter{
		UserID:      userID,
		ServiceName: r.ServiceName,
		PeriodStart: *periodStart,
		PeriodEnd:   *periodEnd,
	}, nil
}

func FromSubscriptionDTO(d dto.SubscriptionDTO) *SubscriptionResponse {
	endDate := FromTime(d.EndDate)
	return &SubscriptionResponse{
		ID:          d.ID.String(),
		UserID:      d.UserID.String(),
		ServiceName: d.ServiceName,
		Price:       d.Price,
		StartDate:   FromTime(&d.StartDate),
		EndDate:     &endDate,
	}
}
