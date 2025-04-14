package provider

import (
	"context"
	"data-provider-service/internal/model"
	repository "data-provider-service/internal/provider/repository"
	"data-provider-service/mapper"
	"fmt"
	api "github.com/5krotov/task-resolver-pkg/api/v1"
	"time"
)

type Provider struct {
	repository repository.ProviderRepository
	mapper     *mapper.Mapper
}

func NewProvider(repository repository.ProviderRepository) *Provider {
	return &Provider{repository: repository, mapper: &mapper.Mapper{}}
}

func (p *Provider) CreateTask(request *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	task := model.Task{
		Name:       request.Name,
		Difficulty: request.Difficulty,
		Status:     0,
		LastUpdate: time.Now(),
	}
	createdTask, err := p.repository.CreateTask(context.Background(), task)
	if err != nil {
		return nil, fmt.Errorf("failed create task: %v", err)
	}
	apiTask := p.mapper.TaskAndStatusesToAPITask(*createdTask, []model.Status{{Status: task.Status, Timestamp: task.LastUpdate}})

	return &api.CreateTaskResponse{Task: apiTask}, nil
}

func (p *Provider) GetTask(taskId int64) (*api.GetTaskResponse, error) {
	task, err := p.repository.GetTaskByID(context.Background(), taskId)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %v", err)
	}
	statuses, err := p.repository.GetStatusesByTaskID(context.Background(), taskId)
	if err != nil {
		return nil, fmt.Errorf("get statuses failed: %v", err)
	}
	apiTask := p.mapper.TaskAndStatusesToAPITask(*task, statuses)

	return &api.GetTaskResponse{Task: apiTask}, nil
}

func (p *Provider) UpdateTaskStatus(request *api.UpdateStatusRequest) error {
	status := model.Status{Status: request.Status.Status, Timestamp: request.Status.Timestamp, TaskID: request.Id}

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
