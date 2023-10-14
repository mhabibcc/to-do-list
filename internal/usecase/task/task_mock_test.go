package task

import (
	"context"
	model "to-do-list/internal/model/task"
)

var TaskRepository = &TaskRepositoryMock{}

type TaskRepositoryMock struct {
	GetAllFunc func(ctx context.Context) ([]model.TaskModel, error)
	CreateFunc func(ctx context.Context, task model.TaskModel) model.TaskModel
	UpdateFunc func(ctx context.Context, task model.TaskModel) (bool, model.TaskModel)
	DeleteFunc func(ctx context.Context, task model.TaskModel) bool
}

func (repository *TaskRepositoryMock) GetAll(ctx context.Context) ([]model.TaskModel, error) {
	return repository.GetAllFunc(ctx)
}

func (repository *TaskRepositoryMock) Create(ctx context.Context, task model.TaskModel) model.TaskModel {
	return repository.CreateFunc(ctx, task)
}

func (repository *TaskRepositoryMock) Update(ctx context.Context, task model.TaskModel) (bool, model.TaskModel) {
	return repository.UpdateFunc(ctx, task)
}

func (repository *TaskRepositoryMock) Delete(ctx context.Context, task model.TaskModel) bool {
	return repository.DeleteFunc(ctx, task)
}
