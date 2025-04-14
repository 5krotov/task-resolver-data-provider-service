package cache_repository

import (
	"context"
)

type CacheRepository interface {
	Cache(ctx context.Context, request interface{}, response interface{}) error
	Load(ctx context.Context, request interface{}) ([]byte, error)
	Close(ctx context.Context)
}
