package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	model "to-do-list/internal/model/task"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	rclient *redis.Client
)

func mockDBAndRedis(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *redis.Client) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// defer db.Close()

	mr, err := miniredis.Run()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rclient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return db, mock, rclient
}

func TestRepo_GetAll(t *testing.T) {

	db, mock, rclient := mockDBAndRedis(t)

	ctx := context.Background()

	type arg struct {
		ctx_arg context.Context
	}

	rows := sqlmock.NewRows([]string{"id", "task_name", "is_done"}).
		AddRow(1, "task 1", true).
		AddRow(2, "task 2", false)

	tests := []struct {
		name    string
		arg     arg
		want    []model.TaskModel
		wantErr error
	}{
		{
			name: "case 1 -> get all data",
			arg:  arg{ctx_arg: ctx},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTaskRepository(db, rclient)
			mock.ExpectQuery(`SELECT id, task_name, is_done FROM tasks`).WillReturnRows(rows)
			result, err := repo.GetAll(tt.arg.ctx_arg)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)

			str, err := rclient.Get(rclient.Context(), "tasks").Result()
			var tasks []model.TaskModel
			assert.NoError(t, err, "error redis")

			err = json.Unmarshal([]byte(str), &tasks)
			assert.NoError(t, err, "error redis")
			assert.Equal(t, tt.want, tasks)

		})
	}

}

func TestRepo_Create(t *testing.T) {

	db, mock, rclient := mockDBAndRedis(t)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}

	tests := []struct {
		name string
		args args
		want model.TaskModel
	}{
		{
			name: "case 1 -> create task data",
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					TaskName: "task 1",
					IsDone:   true,
				},
			},
			want: model.TaskModel{
				ID:       1,
				TaskName: "task 1",
				IsDone:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTaskRepository(db, rclient)
			mock.ExpectQuery(`INSERT INTO tasks (.*) RETURNING id`).WillReturnRows(rows)
			result := repo.Create(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.want, result)
		})
	}

}

func TestRepo_Update(t *testing.T) {

	db, mock, rclient := mockDBAndRedis(t)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}

	tests := []struct {
		name       string
		args       args
		mock       func()
		want       model.TaskModel
		wantStatus bool
	}{
		{
			name: "case 1 -> success update task data",
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					ID:       1,
					TaskName: "task 2",
					IsDone:   true,
				},
			},
			mock: func() {
				mock.ExpectExec(`UPDATE tasks SET task_name=(.*), is_done=(.*) WHERE id=(.*)`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: model.TaskModel{
				ID:       1,
				TaskName: "task 2",
				IsDone:   true,
			},
			wantStatus: true,
		},
		{
			name: "case 2 -> no affected update task data",
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					ID:       1,
					TaskName: "task 2",
					IsDone:   true,
				},
			},
			mock: func() {
				mock.ExpectExec(`UPDATE tasks SET task_name=(.*), is_done=(.*) WHERE id=(.*)`).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			want: model.TaskModel{
				ID:       1,
				TaskName: "task 2",
				IsDone:   true,
			},
			wantStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			repo := NewTaskRepository(db, rclient)
			status, result := repo.Update(tt.args.ctx, tt.args.request)

			_, err := rclient.Get(rclient.Context(), "tasks").Result()
			assert.Equal(t, redis.Nil, err)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantStatus, status)
		})
	}

}

func TestRepo_Delete(t *testing.T) {

	db, mock, rclient := mockDBAndRedis(t)

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		request model.TaskModel
	}

	tests := []struct {
		name string
		args args
		mock func()
		want bool
	}{
		{
			name: "case 1 -> success delete task data",
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					ID:       1,
					TaskName: "task 2",
					IsDone:   true,
				},
			},

			mock: func() {
				mock.ExpectExec(`DELETE FROM tasks WHERE id=(.*)`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: true,
		},
		{
			name: "case 2 -> no affected delete task data",
			args: args{
				ctx: ctx,
				request: model.TaskModel{
					ID:       1,
					TaskName: "task 2",
					IsDone:   true,
				},
			},
			mock: func() {
				mock.ExpectExec(`DELETE FROM tasks WHERE id=(.*)`).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			repo := NewTaskRepository(db, rclient)
			result := repo.Delete(tt.args.ctx, tt.args.request)

			_, err := rclient.Get(rclient.Context(), "tasks").Result()
			assert.Equal(t, redis.Nil, err)
			assert.Equal(t, tt.want, result)
		})
	}

}
