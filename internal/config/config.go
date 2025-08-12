package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env:"ENV" env-default:"local"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type ServerConfig struct {
	Port string     `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	CORS CORSConfig `yaml:"cors"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" env:"DB_DRIVER" env-default:"postgres"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASS" env-required:"true"`
	Database string `yaml:"database" env:"DB_NAME" env-required:"true"`
}

type LoggerConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`
}

type CORSConfig struct {
	AllowOrigins     []string      `yaml:"allow_origins" env:"CORS_ALLOW_ORIGINS" env-default:"*"`
	AllowMethods     []string      `yaml:"allow_methods" env:"CORS_ALLOW_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders     []string      `yaml:"allow_headers" env:"CORS_ALLOW_HEADERS" env-default:"Authorization,Content-Type"`
	ExposeHeaders    []string      `yaml:"expose_headers" env:"CORS_EXPOSE_HEADERS"`
	AllowCredentials bool          `yaml:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS" env-default:"true"`
	MaxAge           time.Duration `yaml:"max_age" env:"CORS_MAX_AGE" env-default:"3600s"`
}

// Examples:
//
//	Postgres: host=localhost user=postgres password=pass dbname=mydb port=5432 sslmode=disable
//	MySQL:    user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local
//	SQLite:   ./app.db or :memory:
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "postgres", "postgresql":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			c.Host, c.Username, c.Password, c.Database, c.Port,
		)

	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database,
		)

	case "sqlite", "sqlite3":
		// path .db or ":memory:"
		return c.Database

	default:
		return ""
	}
}

// Examples:
//
//	Postgres: postgres://postgres:pass@localhost:5432/mydb?sslmode=disable
//	MySQL:    mysql://user:pass@tcp(localhost:3306)/mydb?parseTime=true
//	SQLite:   sqlite3://./app.db  or  sqlite3://:memory:
func (c *DatabaseConfig) GetURL() string {
	switch c.Driver {
	case "postgres", "postgresql":
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			url.QueryEscape(c.Username),
			url.QueryEscape(c.Password),
			c.Host,
			c.Port,
			c.Database,
		)

	case "mysql":
		return fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true",
			url.QueryEscape(c.Username),
			url.QueryEscape(c.Password),
			c.Host,
			c.Port,
			c.Database,
		)

	case "sqlite", "sqlite3":
		if c.Database == ":memory:" {
			return "sqlite3://:memory:"
		}
		return fmt.Sprintf("sqlite3://%s", c.Database)

	default:
		return ""
	}
}

func GetConfig() *Config {
	configPath := fetchConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("config file does not exist: %s", configPath)
	}

	return GetConfigFromPath(configPath)
}

func GetConfigFromPath(path string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Panicf("cannot read config: %s", path)
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	if path == "" {
		path = "configs/config.yaml"
	}

	return path
}
