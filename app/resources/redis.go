package resources

import (
	"github.com/sanches1984/auth/config"
	"github.com/sanches1984/auth/pkg/redis"
)

func InitRedis() (*redis.Client, error) {
	return redis.NewClient(redis.Config{
		Host:              config.Redis().Host,
		Password:          config.Secrets().RedisPassword,
		Db:                config.Redis().Db,
		ConnectionTimeout: config.Redis().ConnectionTimeout,
		OperationTimeout:  config.Redis().OperationTimeout,
	})
}
