package service

import (
	"data-provider-service/internal/cache"
	"data-provider-service/internal/entity"
	"data-provider-service/internal/provider"
	"fmt"
	api "github.com/5krotov/task-resolver-pkg/api/v1"
	"log"
)

type DataProviderService struct {
	cache    *cache.Cache
	provider *provider.Provider
}

func NewDataProviderService(cache *cache.Cache, provider *provider.Provider) *DataProviderService {
	return &DataProviderService{cache: cache, provider: provider}
}

func (s *DataProviderService) SearchTask(request *entity.SearchTaskParams) (*api.SearchTaskResponse, error) {
	return s.provider.SearchTask(request)
}

func (s *DataProviderService) CreateTask(request *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	return s.provider.CreateTask(request)
}

func (s *DataProviderService) GetTask(taskId int64) (*api.GetTaskResponse, error) {
	response, err := s.cache.LoadGetTaskResponse(taskId)
	if err == nil && response != nil {
		log.Println("record getted from cache")
		return response, nil
	}
	if err != nil {
		log.Printf("get record from cache failed: %v \n", err)
	} else {
		log.Println("no record in cache")
	}

	response, err = s.provider.GetTask(taskId)
	if err != nil {
		return nil, fmt.Errorf("get task from provider failed: %v", err)
	}

	err = s.cache.Cache(taskId, response)
	if err != nil {
		log.Printf("cache record failed: %v \n", err)
	}

	return response, nil
}

func (s *DataProviderService) UpdateStatus(request *api.UpdateStatusRequest) error {
	return s.provider.UpdateTaskStatus(request)
}
