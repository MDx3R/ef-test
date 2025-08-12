package gorm

import (
	"fmt"

	"github.com/MDx3R/ef-test/internal/domain/entity"
	gormmodel "github.com/MDx3R/ef-test/internal/infra/database/gorm/model"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormSubscriptionRepository struct {
	tx *gorm.DB
}

func NewGormSubscriptionRepository(db *gorm.DB) usecase.SubscriptionRepository {
	return &gormSubscriptionRepository{db}
}

func (r *gormSubscriptionRepository) Get(id uuid.UUID) (*entity.Subscription, error) {
	var sub gormmodel.SubscriptionModel

	err := r.tx.First(&sub, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, usecase.ErrNotFound
		}
		return nil, wrap(usecase.ErrRepository, err)
	}

	return sub.ToEntity(), nil
}
func (r *gormSubscriptionRepository) List(filter dto.SubscriptionFilter) ([]*entity.Subscription, error) {
	var subs []gormmodel.SubscriptionModel

	stmt := r.tx
	if filter.UserID != nil {
		stmt = stmt.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceName != nil {
		stmt = stmt.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.StartDate != nil {
		stmt = stmt.Where("start_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		stmt = stmt.Where("end_date IS NULL OR end_date <= ?", *filter.EndDate)
	}

	offset := (filter.Page - 1) * filter.PageSize
	err := stmt.Offset(offset).Limit(filter.PageSize).Find(&subs).Error

	if err != nil {
		return nil, wrap(usecase.ErrRepository, err)
	}

	result := make([]*entity.Subscription, len(subs))
	for i, sub := range subs {
		result[i] = sub.ToEntity()
	}

	return result, nil
}
func (r *gormSubscriptionRepository) Add(sub *entity.Subscription) error {
	model := gormmodel.FromEntity(sub)

	err := r.tx.Create(model).Error
	if err != nil {
		return wrap(usecase.ErrRepository, err)
	}
	return nil
}
func (r *gormSubscriptionRepository) Update(sub *entity.Subscription) error {
	model := gormmodel.FromEntity(sub)

	err := r.tx.Save(model).Error
	if err != nil {
		return wrap(usecase.ErrRepository, err)
	}
	return nil
}
func (r *gormSubscriptionRepository) Delete(id uuid.UUID) error {
	err := r.tx.Delete(&gormmodel.SubscriptionModel{}, "id = ?", id).Error
	if err != nil {
		return wrap(usecase.ErrRepository, err)
	}
	return nil
}
func (r *gormSubscriptionRepository) CalculateTotalCost(filter dto.TotalCostFilter) (int, error) {
	var result int

	// Calculate the total subscription cost for the given user and service,
	// considering only subscriptions whose start_date falls within the specified period.
	stmt := r.tx.Model(&gormmodel.SubscriptionModel{}).Select("sum(price)")
	stmt = stmt.Where("user_id = ?", filter.UserID)
	stmt = stmt.Where("service_name = ?", filter.ServiceName)
	stmt = stmt.Where("start_date BETWEEN ? AND ?", filter.PeriodStart, filter.PeriodEnd)

	err := stmt.Scan(&result).Error
	if err != nil {
		return 0, wrap(usecase.ErrRepository, err)
	}
	return result, nil
}

func wrap(to, with error) error {
	return fmt.Errorf("%w: %v", to, with)
}
