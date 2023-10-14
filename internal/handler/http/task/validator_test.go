package task

import (
	"testing"
	model "to-do-list/internal/model/task"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	tests := []struct {
		name string
		arg  model.TaskModel
		want []model.ErrorField
	}{
		{
			name: "case 1 -> task name filled",
			arg: model.TaskModel{
				TaskName: "test1",
			},
			want: nil,
		},
		{
			name: "case 2 -> task name not filled",
			arg:  model.TaskModel{},
			want: []model.ErrorField{
				{
					FieldName: "TaskName",
					Message:   "TaskName is required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.arg)

			assert.Equal(t, tt.want, err)
		})
	}
}
