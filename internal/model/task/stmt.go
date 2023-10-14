package task

const FetchAllTaskQuery = `SELECT id, task_name, is_done FROM tasks`

const InsertTaskReturnIdQuery = `INSERT INTO tasks (task_name, is_done) VALUES ($1, $2) RETURNING id`

const UpdateTaskQuery = `UPDATE tasks SET task_name=$1, is_done=$2 WHERE id=$3`

const DeleteTaskQuery = `DELETE FROM tasks WHERE id=$1`
