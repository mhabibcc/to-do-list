package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	model "to-do-list/internal/model/task"

	"github.com/go-redis/redis/v8"
)

type Repo struct {
	Db    *sql.DB
	Redis *redis.Client
}

const (
	redisTaskGetAll = "tasks"
)

func NewTaskRepository(db *sql.DB, redis *redis.Client) *Repo {
	return &Repo{
		Db:    db,
		Redis: redis,
	}
}

func (r *Repo) GetAll(ctx context.Context) ([]model.TaskModel, error) {

	var Tasks = []model.TaskModel{}
	rdb, err := r.Redis.Get(ctx, redisTaskGetAll).Result()
	if err != redis.Nil {
		fmt.Println(err)
	} else {
		err = json.Unmarshal([]byte(rdb), &Tasks)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			if len(Tasks) > 0 {
				return Tasks, nil
			}
		}
	}

	rows, err := r.Db.QueryContext(ctx, model.FetchAllTaskQuery)

	if err != nil {
		fmt.Println("Error on Repo :", err)
		return nil, errors.New("database error")
	}

	defer rows.Close()

	task_row := model.TaskModel{}

	for rows.Next() {
		err := rows.Scan(&task_row.ID, &task_row.TaskName, &task_row.IsDone)
		if err != nil {
			fmt.Println(err)
		}
		Tasks = append(Tasks, task_row)
	}

	json_data, err := json.Marshal(Tasks)

	if err != nil {
		fmt.Println(err)
	} else {
		_ = r.Redis.Set(context.Background(), redisTaskGetAll, json_data, time.Duration(0))
	}

	return Tasks, nil
}

func (r *Repo) Create(ctx context.Context, task model.TaskModel) model.TaskModel {

	err := r.Db.QueryRowContext(ctx, model.InsertTaskReturnIdQuery, task.TaskName, task.IsDone).Scan(&task.ID)

	if err != nil {
		fmt.Println(err)
	}

	_ = r.Redis.Del(ctx, redisTaskGetAll)

	return task
}

func (r *Repo) Update(ctx context.Context, task model.TaskModel) (bool, model.TaskModel) {

	res, err := r.Db.ExecContext(ctx, model.UpdateTaskQuery, task.TaskName, task.IsDone, task.ID)

	if err != nil {
		fmt.Println(err)
	}

	result := false
	affected, _ := res.RowsAffected()
	if affected > 0 {
		result = true
	}

	_ = r.Redis.Del(context.Background(), redisTaskGetAll)

	return result, task
}

func (r *Repo) Delete(ctx context.Context, task model.TaskModel) bool {

	res, err := r.Db.ExecContext(ctx, model.DeleteTaskQuery, task.ID)

	if err != nil {
		fmt.Println(err)
	}

	result := false
	affected, _ := res.RowsAffected()
	if affected > 0 {
		result = true
	}

	_ = r.Redis.Del(ctx, redisTaskGetAll)

	return result
}
