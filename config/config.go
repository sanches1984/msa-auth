package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"time"
)

// Global config var
var config appConfig
var env string

const App = "auth"

type envConfig struct {
	Addr  string            `yaml:"addr"`
	SQL   SQLConfig         `yaml:"sql"`
	Redis RedisConfig       `yaml:"redis"`
	Token TokenConfig       `yaml:"token"`
	Env   map[string]string `yaml:"env"`
}

type SQLConfig struct {
	Host string `yaml:"host"`
	Db   string `yaml:"db"`
	User string `yaml:"user"`
	Pool int    `yaml:"pool"`
}

type RedisConfig struct {
	Host              string        `yaml:"host"`
	Db                int           `yaml:"db"`
	ConnectionTimeout time.Duration `yaml:"conn_timeout"`
	OperationTimeout  time.Duration `yaml:"oper_timeout"`
}

type TokenConfig struct {
	AccessTTL  time.Duration `yaml:"access_ttl"`
	RefreshTTL time.Duration `yaml:"refresh_ttl"`
}

type SecretsConfig struct {
	SQLPassword   string
	RedisPassword string
	JwtSecret     string
}

type appConfig struct {
	// per env conf
	env     map[string]envConfig
	secrets SecretsConfig
}

func Env() string {
	if env != "" {
		return env
	}

	env, exists := os.LookupEnv("AUTH_ENV")
	if exists != true {
		env = "default"
	}

	return env
}

func Load() error {
	if err := loadYaml(); err != nil {
		return err
	}
	if err := loadEnv(); err != nil {
		return err
	}
	return nil
}

func Addr() string {
	return config.env[Env()].Addr
}

func AccessTokenTTL() time.Duration {
	return config.env[Env()].Token.AccessTTL
}

func RefreshTokenTTL() time.Duration {
	return config.env[Env()].Token.RefreshTTL
}

func Redis() RedisConfig {
	return config.env[Env()].Redis
}

func SQLDSN(noDb ...bool) string {
	if len(noDb) == 1 && noDb[0] == true {
		return "postgres://" + url.QueryEscape(config.env[Env()].SQL.User) + ":" + url.QueryEscape(config.secrets.SQLPassword) +
			"@" + config.env[Env()].SQL.Host + "/postgres?sslmode=disable"
	}

	return "postgres://" + url.QueryEscape(config.env[Env()].SQL.User) + ":" + url.QueryEscape(config.secrets.SQLPassword) +
		"@" + config.env[Env()].SQL.Host + "/" + config.env[Env()].SQL.Db +
		"?sslmode=disable"
}

func Secrets() SecretsConfig {
	return config.secrets
}

func loadYaml() error {
	var err error
	var cfg = io.ReadCloser(os.Stdin)

	cfgPath := getConfigPath("config.yaml")
	if cfg, err = os.Open(cfgPath); err != nil {
		return fmt.Errorf("failed to open config: %w", err)
	}

	data, err := ioutil.ReadAll(cfg)
	_ = cfg.Close()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	err = yaml.Unmarshal(data, &config.env)
	if err != nil {
		return fmt.Errorf("bad yaml format: %w", err)
	}

	return nil
}

func loadEnv() error {
	var err error
	config.secrets.SQLPassword, err = getEnv("AUTH_SQL_PASSWORD")
	if err != nil {
		return err
	}
	config.secrets.RedisPassword, err = getEnv("AUTH_REDIS_PASSWORD")
	if err != nil {
		return err
	}
	config.secrets.JwtSecret, err = getEnv("AUTH_JWT_SECRET")
	if err != nil {
		return err
	}
	return nil
}

func getEnv(v string, defaultValue ...string) (string, error) {
	// from yaml config
	if len(config.env[Env()].Env) > 0 {
		if env, ok := config.env[Env()].Env[v]; ok {
			return env, nil
		}
	}
	// from environment variable
	env, exists := os.LookupEnv(v)
	if exists != true {
		if len(defaultValue) != 0 {
			return defaultValue[0], nil
		}

		return "", fmt.Errorf("unable to find env var: %s", v)
	}

	return env, nil
}

func getRootPath() string {
	dirList := [...]string{
		"/app",
		path.Join(os.Getenv("GOPATH"), "src/github.com/sanches1984/auth"),
		".",
	}
	for _, dir := range dirList {
		if _, err := os.Stat(dir + "/config"); !os.IsNotExist(err) {
			return dir
		}
	}
	panic("root path not found")
}

func getConfigPath(name string) string {
	return path.Join(getRootPath(), "config", name)
}
