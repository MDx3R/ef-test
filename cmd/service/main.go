package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/MDx3R/ef-test/internal/infra/app"
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

	app := app.NewApp(cfg, logger)

	logger.Info("running migrations...")

	if err := app.Database.Migrate(); err != nil {
		logger.Fatalf("failed to run migrations: %v", err)
	}

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
