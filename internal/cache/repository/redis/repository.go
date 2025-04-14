package redis

import (
	"context"
	"data-provider-service/internal/config"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Repository struct {
	client        *redis.Client
	cacheLifetime time.Duration
}

func NewRepository(cfg config.RedisConfig) (*Repository, error) {
	password, exists := os.LookupEnv(cfg.PasswordEnvVar)
	if !exists {
		return nil, fmt.Errorf("password env var not found")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DataBase,
		Password: password,
	})

	cacheLifeTime, err := time.ParseDuration(cfg.CacheLifetime)
	if err != nil {
		return nil, fmt.Errorf("parsing lifetime failed: %v", err)
	}
	return &Repository{client: client, cacheLifetime: cacheLifeTime}, nil
}

func (p *Repository) Cache(ctx context.Context, request interface{}, response interface{}) error {
	key, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	value, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if err := p.client.Set(ctx, string(key), value, p.cacheLifetime).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

func (p *Repository) Load(ctx context.Context, request interface{}) ([]byte, error) {

	key, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	data, err := p.client.Get(ctx, string(key)).Bytes()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	return data, nil
}

func (p *Repository) Close(ctx context.Context) {
	p.client.Shutdown(ctx)
}
