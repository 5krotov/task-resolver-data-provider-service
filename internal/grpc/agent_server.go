package grpc

import (
	"context"
	"data-provider-service/internal/service"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataProviderServiceAgent struct {
	pb.UnimplementedDataProviderServiceServer
	service *service.DataProviderService
}

func NewDataProviderServiceAgent(service *service.DataProviderService) *DataProviderServiceAgent {
	return &DataProviderServiceAgent{service: service}
}

func (s *DataProviderServiceAgent) SearchTask(ctx context.Context, request *pb.SearchTaskRequest) (*pb.SearchTaskResponse, error) {
	return s.service.SearchTask(request)
}

func (s *DataProviderServiceAgent) CreateTask(ctx context.Context, request *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	return s.service.CreateTask(request)
}

func (s *DataProviderServiceAgent) GetTask(ctx context.Context, request *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	return s.service.GetTask(request)
}

func (s *DataProviderServiceAgent) UpdateStatus(ctx context.Context, request *pb.UpdateStatusRequest) (*emptypb.Empty, error) {
	return s.service.UpdateStatus(request)
}
