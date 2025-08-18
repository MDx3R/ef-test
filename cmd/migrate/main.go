package main

import (
	"log"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/MDx3R/ef-test/internal/infra/database/migrate"
	logruslogger "github.com/MDx3R/ef-test/internal/infra/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	cfg := config.GetConfig()
	logger := logruslogger.NewLogger()
	logger = logruslogger.SetupLogger(logger, cfg)

	logruslogger.LogConfig(logger, cfg)

	logger.Info("running migrations...")
	migrate.MustRunMigrations(&cfg.Database, logger)
	logger.Info("finished migrations...")
}
