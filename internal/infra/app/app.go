package app

import (
	"context"
	"time"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/MDx3R/ef-test/internal/infra/database/gorm"
	ginserver "github.com/MDx3R/ef-test/internal/infra/server/gin"
	ginware "github.com/MDx3R/ef-test/internal/infra/server/gin/middleware"
	handlers "github.com/MDx3R/ef-test/internal/transport/http/gin"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

type App struct {
	Config   *config.Config
	Server   *ginserver.GinServer
	Database *gorm.GormDatabase
	Logger   *logrus.Logger
}

func NewApp(cfg *config.Config, logger *logrus.Logger) *App {
	logger.Info("establishing database connection")

	gormDB, err := gorm.NewGormDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("failed to create database: %v", err)
	}

	logger.Info("database connected")

	subRepository := gorm.NewGormSubscriptionRepository(gormDB.GetDB())

	subService := usecase.NewSubscriptionService(subRepository)

	subHandler := handlers.NewSubscriptionHandler(subService, logger)

	logger.Info("initializing http server")

	ginserver.SetMode(cfg)
	server := ginserver.New(&cfg.Server)

	server.UseMiddleware(
		ginware.NewCORSMiddleware(&cfg.Server.CORS),
		gin.Recovery(),
		ginware.LoggerMiddleware(logger),
	)

	server.RegisterSwagger()
	server.RegisterSubscriptionHandler(subHandler)

	logger.Info("http server initialized")

	return &App{Config: cfg, Server: server, Database: gormDB, Logger: logger}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.Logger.Fatalf("server failed to run: %v", err)
	}
}

func (a *App) Run() error {
	a.Logger.Infof("starting server on port %s", a.Config.Server.Port)
	if err := a.Server.Run(); err != nil {
		a.Logger.Errorf("server failed to start: %v", err)
		return err
	}
	return nil
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.Logger.Info("shutting down server...")
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Errorf("failed to shutdown server: %v", err)
	}

	a.Logger.Info("closing database connection...")
	if err := a.Database.Dispose(); err != nil {
		a.Logger.Errorf("failed to shutdown database: %v", err)
	}

	a.Logger.Info("shutdown complete")
	return nil
}
