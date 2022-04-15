package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sanches1984/gopkg-logger"
	"time"
)

// Global config var
var config appConfig

type appConfig struct {
	Env Environment
}

type Environment struct {
	AppName        string        `envconfig:"APP_NAME"        default:"auth"`
	Host           string        `envconfig:"HOST"            required:"true"`
	SQLDSN         string        `envconfig:"SQLDSN"          required:"true"`
	MigrationsPath string        `envconfig:"MIGRATIONS_PATH" default:"internal/pkg/migrations"`
	RedisHost      string        `envconfig:"REDIS_HOST"      required:"true"`
	RedisPassword  string        `envconfig:"REDIS_PASSWORD"`
	JwtSecret      string        `envconfig:"JWT_SECRET"      required:"true"`
	ConnectTimeout time.Duration `envconfig:"CONNECT_TIMEOUT" default:"5s"`
	ReadTimeout    time.Duration `envconfig:"READ_TIMEOUT"    default:"2s"`
	AccessTTL      time.Duration `envconfig:"ACCESS_TTL"      default:"6h"`
	RefreshTTL     time.Duration `envconfig:"REFRESH_TTL"     default:"24h"`
	MetricsHost    string        `envconfig:"METRICS_HOST"    default:"localhost:8080"`
	LogType        log.Type      `envconfig:"LOG_TYPE"        default:"console"`
	LogLevel       log.Level     `envconfig:"LOG_LEVEL"       default:"info"`
}

func Load() error {
	if err := godotenv.Load(); err != nil {
		return errors.New(".env file not found")
	}
	return envconfig.Process("AUTH", &config.Env)
}

func Env() Environment {
	return config.Env
}
