package task

const ListTableName = "task"

// swagger:model Task 
type TaskModel struct {
	// ID of task
	// in: int64
	ID       int64  `json:"id"`
	// Name of task
	// in: string
	TaskName string `json:"task_name" validate:"required"`
	// Task status
	// in: bool
	IsDone   bool   `json:"is_done"`
}

type ValidationResponse struct {
	Message string       `json:"message"`
	Error   []ErrorField `json:"error"`
}

type ErrorField struct {
	FieldName string `json:"field"`
	Message   string `json:"message"`
}
