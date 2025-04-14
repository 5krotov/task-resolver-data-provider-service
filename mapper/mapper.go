package mapper

import (
	"data-provider-service/internal/model"
	entity "github.com/5krotov/task-resolver-pkg/entity/v1"
)

type Mapper struct {
}

func (m *Mapper) StatusToAPIStatus(status model.Status) entity.Status {
	return entity.Status{
		Status:    status.Status,
		Timestamp: status.Timestamp,
	}
}

func (m *Mapper) TaskAndStatusesToAPITask(task model.Task, statuses []model.Status) entity.Task {
	var apiStatuses []entity.Status

	for _, status := range statuses {
		apiStatuses = append(apiStatuses, m.StatusToAPIStatus(status))
	}

	return entity.Task{
		Id:            task.ID,
		Name:          task.Name,
		Difficulty:    task.Difficulty,
		StatusHistory: apiStatuses,
	}
}
