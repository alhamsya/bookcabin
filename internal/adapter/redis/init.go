package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/alhamsya/bookcabin/internal/core/domain/constant"
	"github.com/alhamsya/bookcabin/internal/core/port/repository"
	"github.com/alhamsya/bookcabin/pkg/manager/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context, cfg *config.Application) (port.CacheRepo, error) {
	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Static.Redis.Host, cfg.Static.Redis.Port),
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(constant.NetworkTCP, addr)
		},
		Username: cfg.Credential.Redis.Username,
		Password: cfg.Credential.Redis.Password,
		DB:       cfg.Static.Redis.DB,
	}

	client := redis.NewClient(opts)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("ping result: %w", err)
	}

	return &Redis{client}, nil
}
