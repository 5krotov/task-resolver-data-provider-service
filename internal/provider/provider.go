package provider

import (
	"context"
	"data-provider-service/internal/entity"
	"data-provider-service/internal/mapper"
	"data-provider-service/internal/model"
	repository "data-provider-service/internal/provider/repository"
	"fmt"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"math"
	"time"
)

type Provider struct {
	repository repository.ProviderRepository
	mapper     *mapper.Mapper
}

func NewProvider(repository repository.ProviderRepository) *Provider {
	return &Provider{repository: repository, mapper: &mapper.Mapper{}}
}

func (p *Provider) SearchTask(params *pb.SearchTaskRequest) (*pb.SearchTaskResponse, error) {
	count, err := p.repository.CountTask(context.Background())
	if err != nil {
		return nil, fmt.Errorf("count tasks failed: %v", err)
	}
	countPages := count / int64(params.PerPage)
	if count%int64(params.PerPage) != 0 {
		countPages++
	}
	var countPagesInt int
	if countPages > int64(math.MaxInt) {
		countPagesInt = math.MaxInt
	} else {
		countPagesInt = int(countPages)
	}
	foundTasks, foundTasksStatuses, err := p.repository.SearchTask(context.Background(), &entity.SearchTaskParams{int(params.PerPage), int(params.Page)})
	if err != nil {
		return nil, fmt.Errorf("finding tasks failed: %v", err)
	}

	apiTasks := make([]*pb.Task, len(foundTasks))

	for ind, task := range foundTasks {
		apiTask := p.mapper.TaskAndStatusesToAPITask(task, foundTasksStatuses[ind])
		apiTasks[ind] = &apiTask
	}

	return &pb.SearchTaskResponse{Pages: int64(countPagesInt), Tasks: apiTasks}, nil
}

func (p *Provider) CreateTask(request *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task := model.Task{
		Name:       request.Name,
		Difficulty: int(request.Difficulty),
		Status:     0,
		LastUpdate: time.Now(),
	}
	createdTask, err := p.repository.CreateTask(context.Background(), task)
	if err != nil {
		return nil, fmt.Errorf("failed create task: %v", err)
	}
	apiTask := p.mapper.TaskAndStatusesToAPITask(*createdTask, []model.Status{{Status: task.Status, Timestamp: task.LastUpdate}})

	return &pb.CreateTaskResponse{Task: &apiTask}, nil
}

func (p *Provider) GetTask(taskId int64) (*pb.GetTaskResponse, error) {
	task, err := p.repository.GetTaskByID(context.Background(), taskId)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %v", err)
	}
	statuses, err := p.repository.GetStatusesByTaskID(context.Background(), taskId)
	if err != nil {
		return nil, fmt.Errorf("get statuses failed: %v", err)
	}
	apiTask := p.mapper.TaskAndStatusesToAPITask(*task, statuses)

	return &pb.GetTaskResponse{Task: &apiTask}, nil
}

func (p *Provider) UpdateTaskStatus(request *pb.UpdateStatusRequest) error {
	status := model.Status{Status: int(request.Status.Status), Timestamp: request.Status.Timestamp.AsTime(), TaskID: request.Id}

	err := p.repository.UpdateTaskStatus(context.Background(), status)
	if err != nil {
		return fmt.Errorf("update status failed: %v", err)
	}

	return nil
}

func (p *Provider) Stop() error {
	err := p.repository.Close()
	if err != nil {
		return fmt.Errorf("close repository failed: %v", err)
	}

	return nil
}
