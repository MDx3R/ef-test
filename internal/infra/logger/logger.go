package logger

import (
	"os"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	return logrus.New()
}

func SetupLogger(logger *logrus.Logger, cfg *config.Config) *logrus.Logger {
	logger.SetOutput(os.Stdout)

	setupLevel(logger, cfg)
	setupFormat(logger, cfg)
	setupReportCaller(logger, cfg)

	return logger
}

func setupLevel(logger *logrus.Logger, cfg *config.Config) {
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
}

func setupFormat(logger *logrus.Logger, cfg *config.Config) {
	switch cfg.Logger.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: cfg.Env == "local",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   cfg.Env == "local",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

func setupReportCaller(logger *logrus.Logger, cfg *config.Config) {
	if logger.Level <= logrus.DebugLevel && (cfg.Env == "debug" || cfg.Env == "local") {
		logger.SetReportCaller(true)
	} else {
		logger.SetReportCaller(false)
	}
}

func LogConfig(logger *logrus.Logger, cfg *config.Config) {
	logger.WithFields(logrus.Fields{
		"env": cfg.Env,
		"server": logrus.Fields{
			"port": cfg.Server.Port,
		},
		"database": logrus.Fields{
			"driver":   cfg.Database.Driver,
			"username": cfg.Database.Username,
			"host":     cfg.Database.Host,
			"port":     cfg.Database.Port,
			"database": cfg.Database.Database,
		},
		"logger": logrus.Fields{
			"level":  cfg.Logger.Level,
			"format": cfg.Logger.Format,
		},
	}).Info("loaded configuration")
}
