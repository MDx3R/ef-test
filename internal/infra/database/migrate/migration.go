package migrate

import (
	"fmt"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func MustRunMigrations(cfg *config.DatabaseConfig, logger *logrus.Logger) {
	if err := RunMigrations(cfg, logger); err != nil {
		logger.Fatalf("failed to run migrations: %v", err)
	}
}

func RunMigrations(cfg *config.DatabaseConfig, logger *logrus.Logger) error {
	url := cfg.GetURL()

	m, err := migrate.New(
		"file://./migrations",
		url,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("no migration needed, schema is up to date")
		} else {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	} else {
		logger.Info("migrations applied successfully")
	}
	return nil
}
