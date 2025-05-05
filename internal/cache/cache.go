package cache

import (
	"context"
	repository "data-provider-service/internal/cache/repository"
	"encoding/json"
	"fmt"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
)

type Cache struct {
	repository repository.CacheRepository
}

func NewCache(repository repository.CacheRepository) *Cache {
	return &Cache{repository: repository}
}

func (c *Cache) Cache(request interface{}, response interface{}) error {
	err := c.repository.Cache(context.Background(), request, response)
	if err != nil {
		return fmt.Errorf("cache failed: %v", err)
	}
	return nil
}

func (c *Cache) LoadGetTaskResponse(taskId int64) (*pb.GetTaskResponse, error) {
	responseBytes, err := c.repository.Load(context.Background(), taskId)
	if err != nil {
		return nil, err
	}
	if responseBytes == nil {
		return nil, nil
	}

	var response pb.GetTaskResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	return &response, nil
}
