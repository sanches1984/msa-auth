package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
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
	RedisHost      string        `envconfig:"REDIS_HOST"      required:"true"`
	RedisPassword  string        `envconfig:"REDIS_PASSWORD"`
	JwtSecret      string        `envconfig:"JWT_SECRET"      required:"true"`
	ConnectTimeout time.Duration `envconfig:"CONNECT_TIMEOUT" default:"5s"`
	ReadTimeout    time.Duration `envconfig:"READ_TIMEOUT"    default:"2s"`
	AccessTTL      time.Duration `envconfig:"ACCESS_TTL"      default:"6h"`
	RefreshTTL     time.Duration `envconfig:"REFRESH_TTL"     default:"24h"`
	LogType        string        `envconfig:"LOG"             default:"console"`
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
