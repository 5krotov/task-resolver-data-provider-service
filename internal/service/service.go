package service

import (
	"data-provider-service/internal/cache"
	"data-provider-service/internal/provider"
	"fmt"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type DataProviderService struct {
	cache    *cache.Cache
	provider *provider.Provider
}

func NewDataProviderService(cache *cache.Cache, provider *provider.Provider) *DataProviderService {
	return &DataProviderService{cache: cache, provider: provider}
}

func (s *DataProviderService) SearchTask(request *pb.SearchTaskRequest) (*pb.SearchTaskResponse, error) {
	return s.provider.SearchTask(request)
}

func (s *DataProviderService) CreateTask(request *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	return s.provider.CreateTask(request)
}

func (s *DataProviderService) GetTask(request *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	response, err := s.cache.LoadGetTaskResponse(request.GetId())
	if err == nil && response != nil {
		log.Println("record getted from cache")
		return response, nil
	}
	if err != nil {
		log.Printf("get record from cache failed: %v \n", err)
	} else {
		log.Println("no record in cache")
	}

	response, err = s.provider.GetTask(request.GetId())
	if err != nil {
		return nil, fmt.Errorf("get task from provider failed: %v", err)
	}

	err = s.cache.Cache(request.GetId(), response)
	if err != nil {
		log.Printf("cache record failed: %v \n", err)
	}

	return response, nil
}

func (s *DataProviderService) UpdateStatus(request *pb.UpdateStatusRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.provider.UpdateTaskStatus(request)
}
