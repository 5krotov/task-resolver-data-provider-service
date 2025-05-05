package mapper

import (
	"data-provider-service/internal/model"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Mapper struct {
}

func (m *Mapper) StatusToAPIStatus(status model.Status) *pb.Status {
	return &pb.Status{
		Status:    pb.StatusValue(status.Status),
		Timestamp: timestamppb.New(status.Timestamp),
	}
}

func (m *Mapper) TaskAndStatusesToAPITask(task model.Task, statuses []model.Status) pb.Task {
	var apiStatuses []*pb.Status

	for _, status := range statuses {
		apiStatuses = append(apiStatuses, m.StatusToAPIStatus(status))
	}

	return pb.Task{
		Id:            task.ID,
		Name:          task.Name,
		Difficulty:    pb.Difficulty(task.Difficulty),
		StatusHistory: apiStatuses,
	}
}
