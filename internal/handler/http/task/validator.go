package task

import (
	"fmt"
	model "to-do-list/internal/model/task"

	"gopkg.in/go-playground/validator.v9"
)

func Validate(request model.TaskModel) []model.ErrorField {
	validate := validator.New()
	err := validate.Struct(request)

	if err != nil {
		var (
			arrErrorField = []model.ErrorField{}
			errorField    = model.ErrorField{}
		)
		for _, err := range err.(validator.ValidationErrors) {
			errorField.FieldName = err.Field()
			errorField.Message = fmt.Sprintf("%v is %v", err.Field(), err.ActualTag())
			arrErrorField = append(arrErrorField, errorField)
		}

		return arrErrorField
	}
	return nil
}
