package gorm

import (
	"fmt"

	"github.com/MDx3R/ef-test/internal/config"
	gormmodel "github.com/MDx3R/ef-test/internal/infra/database/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDatabase struct {
	db  *gorm.DB
	cfg *config.DatabaseConfig
}

func NewGormDatabase(cfg *config.DatabaseConfig) (*GormDatabase, error) {
	db, err := createEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gorm DB: %v", err)
	}
	return &GormDatabase{db, cfg}, nil
}

func (d *GormDatabase) Dispose() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to obtain DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}
	return nil
}

func (d *GormDatabase) Migrate() error {
	err := d.db.AutoMigrate(&gormmodel.SubscriptionModel{})
	if err != nil {
		return fmt.Errorf("failed to migrate DB: %w", err)
	}
	return nil
}

func (d *GormDatabase) GetDB() *gorm.DB {
	return d.db
}

func createEngine(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	return db, nil
}
