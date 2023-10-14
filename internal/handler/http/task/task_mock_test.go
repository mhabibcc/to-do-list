package task

import (
	"context"
	model "to-do-list/internal/model/task"
)

var TaskUseCase = &TaskUsecaseMock{}

type TaskUsecaseMock struct {
	GetAllTaskFunc func(ctx context.Context) ([]model.TaskModel, error)
	CreateTaskFunc func(ctx context.Context, r model.TaskModel) (model.TaskModel, error)
	UpdateTaskFunc func(ctx context.Context, r model.TaskModel) (model.TaskModel, error)
	DeleteTaskFunc func(ctx context.Context, r model.TaskModel) error
}

func (mock *TaskUsecaseMock) GetAllTask(ctx context.Context) ([]model.TaskModel, error) {
	return mock.GetAllTaskFunc(ctx)
}

func (mock *TaskUsecaseMock) CreateTask(ctx context.Context, task model.TaskModel) (model.TaskModel, error) {
	return mock.CreateTaskFunc(ctx, task)
}

func (mock *TaskUsecaseMock) UpdateTask(ctx context.Context, task model.TaskModel) (model.TaskModel, error) {
	return mock.UpdateTaskFunc(ctx, task)
}

func (mock *TaskUsecaseMock) DeleteTask(ctx context.Context, task model.TaskModel) error {
	return mock.DeleteTaskFunc(ctx, task)
}
