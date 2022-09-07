package redis

import (
	"context"
	"encoding/json"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

var Module = fx.Provide(NewRedisCache)

type Params struct {
	fx.In
	Config config.Config
	Log    logger.Logger
}

type redisCache struct {
	client *redis.Client
	log    logger.Logger
}

func NewRedisCache(params Params) (storage.Cacher, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     params.Config.GetString(config.RedisHostKey),
		Password: params.Config.GetString(config.RedisPasswordKey),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "could not connect to redis")
	}

	return &redisCache{
		client: redisClient,
		log:    params.Log,
	}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			c.log.Warn("key not found", zap.Any("key", key))
			return "", models.ErrNotFound
		}
		c.log.Errorf("error while getting key", zap.Any("key", key), zap.Error(err))
		return "", err
	}

	c.log.Debug("found value for given key", zap.Any("key", key))

	return value, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	c.log.Debug("set value", zap.Any("key", key), zap.Any("value", value))
	return c.client.Set(ctx, key, value, duration).Err()
}

func (c *redisCache) GetObj(ctx context.Context, key string, value any) error {
	data, err := c.client.Get(ctx, key).Bytes()

	if err != nil {
		if err == redis.Nil {
			c.log.Warn("key not found", zap.Any("key", key))
			return models.ErrNotFound
		}

		c.log.Errorf("error while getting key", zap.Any("key", key), zap.Error(err))
		return err
	}

	if err = json.Unmarshal(data, value); err != nil {
		c.log.Errorf("error while unmarshalling object", zap.Any("key", key), zap.Error(err))
		return err
	}

	c.log.Info("found value for given key", zap.Any("key", key))
	return nil
}

func (c *redisCache) SetObj(ctx context.Context, key string, value any, duration time.Duration) error {
	c.log.Debug("set object", zap.Any("key", key), zap.Any("value", value))

	data, err := json.Marshal(value)
	if err != nil {
		c.log.Error("could not marshal object", zap.Error(err))
		return err
	}

	if err = c.client.Set(ctx, key, data, duration).Err(); err != nil {
		c.log.Error("could not set object", zap.Error(err))
		return err
	}

	return nil
}

func (c *redisCache) Delete(ctx context.Context, keys ...string) error {
	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		c.log.Error("could not delete keys", zap.Error(err))
		return err
	}

	return nil
}

func (c *redisCache) Close() error {
	return c.client.Close()
}
