package task

import (
	"context"
	"errors"
	"reflect"
	"testing"
	model "to-do-list/internal/model/task"

	"github.com/stretchr/testify/assert"
)

func TestNewUsecase(t *testing.T) {
	type args struct {
		taskRepository *TaskRepositoryMock
	}

	tests := []struct {
		name string
		args args
		want *Usecase
	}{
		{
			name: "case 1 -> success when initiative new task usecase",
			args: args{
				taskRepository: &TaskRepositoryMock{},
			},
			want: &Usecase{
				taskRepo: &TaskRepositoryMock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUseCase(tt.args.taskRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetAllTask(t *testing.T) {
	type fields struct {
		taskRepository *TaskRepositoryMock
	}

	ctx := context.Background()

	type arg struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		arg     arg
		fields  fields
		want    []model.TaskModel
		wantErr error
	}{
		{
			name: "case 1 -> success return data with task model struct",
			arg:  arg{ctx: ctx},
			fields: fields{taskRepository: &TaskRepositoryMock{
				GetAllFunc: func(ctx context.Context) ([]model.TaskModel, error) {
					return []model.TaskModel{
						{
							ID:       1,
							TaskName: "task 1",
							IsDone:   true,
						},
						{
							ID:       2,
							TaskName: "task 2",
							IsDone:   false,
						},
					}, nil
				},
			}},
			want: []model.TaskModel{
				{
					ID:       1,
					TaskName: "task 1",
					IsDone:   true,
				},
				{
					ID:       2,
					TaskName: "task 2",
					IsDone:   false,
				},
			},
			wantErr: nil,
		},
		{
			name: "case 1 -> failed return data",
			arg:  arg{ctx: ctx},
			fields: fields{taskRepository: &TaskRepositoryMock{
				GetAllFunc: func(ctx context.Context) ([]model.TaskModel, error) {
					return nil, errors.New("database error")
				},
			}},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				taskRepo: tt.fields.taskRepository,
			}

			got, err := u.GetAllTask(tt.arg.ctx)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUseCase_CreateTask(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		taskRepository *TaskRepositoryMock
	}

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want1  model.TaskModel
		want2  error
	}{
		{
			name: "case 1 -> success create task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				CreateFunc: func(ctx context.Context, task model.TaskModel) model.TaskModel {
					task.ID = 1
					return task
				},
			}},
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					TaskName: "task 1",
					IsDone:   true,
				},
			},
			want1: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			},
			want2: nil,
		},
		{
			name: "case 2 -> fail create task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				CreateFunc: func(ctx context.Context, task model.TaskModel) model.TaskModel {
					task.ID = 0
					return task
				},
			}},
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					TaskName: "task 1",
					IsDone:   true,
				},
			},
			want1: model.TaskModel{
				ID:       0,
				TaskName: "task 1",
				IsDone:   true,
			},
			want2: errors.New("error, create task failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				taskRepo: tt.fields.taskRepository,
			}

			got, err := u.CreateTask(tt.args.ctx, tt.args.request)

			assert.Equal(t, tt.want1, got)
			assert.Equal(t, tt.want2, err)
		})
	}
}

func TestUseCase_UpdateTask(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		taskRepository *TaskRepositoryMock
	}

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want1  model.TaskModel
		want2  error
	}{
		{
			name: "case 1 -> success update task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				UpdateFunc: func(ctx context.Context, task model.TaskModel) (bool, model.TaskModel) {
					return true, task
				},
			}},
			args: args{ctx: ctx, request: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			}},
			want1: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			},
			want2: nil,
		},
		{
			name: "case 1 -> fail update task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				UpdateFunc: func(ctx context.Context, task model.TaskModel) (bool, model.TaskModel) {
					return false, task
				},
			}},
			args: args{ctx: ctx, request: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			}},
			want1: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			},
			want2: errors.New("error, update task failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				taskRepo: tt.fields.taskRepository,
			}

			got, err := u.UpdateTask(tt.args.ctx, tt.args.request)

			assert.Equal(t, tt.want1, got)
			assert.Equal(t, tt.want2, err)
		})
	}
}

func TestUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		taskRepository *TaskRepositoryMock
	}

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "case 1 -> success update task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				DeleteFunc: func(ctx context.Context, task model.TaskModel) bool {
					return true
				},
			}},
			args: args{ctx: ctx, request: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			}},
			want: nil,
		},
		{
			name: "case 1 -> fail update task",
			fields: fields{taskRepository: &TaskRepositoryMock{
				DeleteFunc: func(ctx context.Context, task model.TaskModel) bool {
					return false
				},
			}},
			args: args{ctx: ctx, request: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			}},
			want: errors.New("error, delete task failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				taskRepo: tt.fields.taskRepository,
			}

			got := u.DeleteTask(tt.args.ctx, tt.args.request)

			assert.Equal(t, tt.want, got)
		})
	}
}
