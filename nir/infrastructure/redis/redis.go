package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/newRational/soc/internal/model"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, resp *model.Response) error
	Get(ctx context.Context, key string) (*model.Response, error)
	Del(ctx context.Context, key string) error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(opt *redis.Options) *Redis {
	return &Redis{
		redis.NewClient(opt),
	}
}

func (r *Redis) Set(ctx context.Context, key string, resp *model.Response) error {
	val, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, val, time.Minute*10).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (*model.Response, error) {
	resp := model.Response{}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &resp)

	return &resp, err
}

func (r *Redis) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()

	return err
}
