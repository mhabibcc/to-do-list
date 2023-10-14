package task

import (
	"context"
	"errors"
	model "to-do-list/internal/model/task"
)

type Usecase struct {
	taskRepo Repo
}

func NewUseCase(repo Repo) *Usecase {
	return &Usecase{
		taskRepo: repo,
	}
}

type Repo interface {
	GetAll(ctx context.Context) ([]model.TaskModel, error)
	Create(ctx context.Context, task model.TaskModel) model.TaskModel
	Update(ctx context.Context, task model.TaskModel) (bool, model.TaskModel)
	Delete(ctx context.Context, task model.TaskModel) bool
}

func (u *Usecase) GetAllTask(ctx context.Context) ([]model.TaskModel, error) {
	tasks, err := u.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *Usecase) CreateTask(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
	task_create := u.taskRepo.Create(ctx, r)
	if task_create.ID == 0 {
		return task_create, errors.New("error, create task failed")
	}
	return task_create, nil
}

func (u *Usecase) UpdateTask(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
	success, task_update := u.taskRepo.Update(ctx, r)
	if !success {
		return task_update, errors.New("error, update task failed")
	}
	return task_update, nil
}

func (u *Usecase) DeleteTask(ctx context.Context, r model.TaskModel) error {
	task_delete := u.taskRepo.Delete(ctx, r)
	if !task_delete {
		return errors.New("error, delete task failed")
	}
	return nil
}
