package resources

import (
	"github.com/rs/zerolog"
	"github.com/sanches1984/msa-auth/config"
	"github.com/sanches1984/msa-auth/pkg/redis"
)

func InitRedis(logger zerolog.Logger) (*redis.Client, error) {
	client, err := redis.NewClient(redis.Config{
		Host:              config.Env().RedisHost,
		Password:          config.Env().RedisPassword,
		ConnectionTimeout: config.Env().ConnectTimeout,
		OperationTimeout:  config.Env().ReadTimeout,
	})
	if err != nil {
		return nil, err
	}

	logger.Info().Str("host", config.Env().RedisHost).Msg("redis connected")
	return client, nil
}
