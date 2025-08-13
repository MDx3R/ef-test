package gin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MDx3R/ef-test/internal/config"
	ginhandlers "github.com/MDx3R/ef-test/internal/transport/http/gin"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	cfg    *config.ServerConfig
	engine *gin.Engine
	server *http.Server
}

func New(cfg *config.ServerConfig) *GinServer {
	engine := gin.New()

	s := &GinServer{
		cfg:    cfg,
		engine: engine,
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		Handler:      s.engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s
}

func SetMode(cfg *config.Config) {
	switch cfg.Env {
	case "test":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// UseMiddleware registers global middleware in Gin.
// Middleware are executed in the order they are registered.
//
// Example execution order:
//  1. Recovery()        // runs before all subsequent middleware
//  2. LoggerMiddleware  // runs after Recovery but before the handler
//  3. Other middleware
//
// It is strongly recommended to register middleware before defining routes,
// so that they apply to all endpoints.
func (g *GinServer) UseMiddleware(mw ...gin.HandlerFunc) {
	g.engine.Use(mw...)
}

func (g *GinServer) RegisterSubscriptionHandler(handler *ginhandlers.SubscriptionHandler) {
	subGroup := g.engine.Group("/subscription")

	subGroup.GET("", handler.List)
	subGroup.GET("/:id", handler.Get)
	subGroup.POST("", handler.Create)
	subGroup.PUT("/:id", handler.Update)
	subGroup.DELETE("/:id", handler.Delete)
	subGroup.GET("/total", handler.CalculateTotalCost)
}

func (g *GinServer) Run() error {
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to run gin server: %w", err)
	}
	return nil
}

func (g *GinServer) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}
