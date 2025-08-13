package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/MDx3R/ef-test/internal/infra/app"
	"github.com/MDx3R/ef-test/internal/infra/database/migrate"
	logruslogger "github.com/MDx3R/ef-test/internal/infra/logger"
	"github.com/joho/godotenv"
)

// @title Effective Mobile GO - Subscription Service API
// @version 1.0
// @description REST API для хранения и обработки информации об онлайн-подписках пользователей.
// @description Сервис позволяет добавлять, изменять, удалять и просматривать записи о подписках, а также рассчитывать суммарную стоимость подписок за выбранный период.

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @tag.name subscriptions
// @tag.description Операции с подписками пользователей

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	cfg := config.GetConfig()
	logger := logruslogger.NewLogger()
	logger = logruslogger.SetupLogger(logger, cfg)

	logruslogger.LogConfig(logger, cfg)

	app := app.NewApp(cfg, logger)

	logger.Info("running migrations...")
	migrate.MustRunMigrations(&cfg.Database, logger)
	logger.Info("finished migrations...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		app.MustRun()
	}()

	<-stop

	logger.Info("shutting down gracefully...")

	if err := app.Shutdown(); err != nil {
		logger.Errorf("hhutdown error: %v", err)
	}
	logger.Info("application stopped gracefully")
}
