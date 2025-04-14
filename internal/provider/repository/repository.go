package provider_repository

import (
	"context"
	"data-provider-service/internal/model"
)

type ProviderRepository interface {
	CreateTask(ctx context.Context, task model.Task) (*model.Task, error)
	UpdateTaskStatus(ctx context.Context, status model.Status) error
	GetTaskByID(ctx context.Context, id int64) (*model.Task, error)
	GetStatusesByTaskID(ctx context.Context, taskID int64) ([]model.Status, error)
	Close() error
}
