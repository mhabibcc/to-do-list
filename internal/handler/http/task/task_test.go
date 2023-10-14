package task

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	model "to-do-list/internal/model/task"
	util "to-do-list/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		taskUseCase *TaskUsecaseMock
	}

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "case 1 -> success when initiative new task handler",
			args: args{
				taskUseCase: &TaskUsecaseMock{},
			},
			want: &Handler{
				useCase: &TaskUsecaseMock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.taskUseCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	type fields struct {
		taskUseCase *TaskUsecaseMock
	}

	tests := []struct {
		name         string
		fields       fields
		wantCode     int
		wantResponse interface{}
	}{
		{
			name: "case 1 -> success when get all task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				GetAllTaskFunc: func(ctx context.Context) ([]model.TaskModel, error) {
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
			wantCode: http.StatusOK,
			wantResponse: []model.TaskModel{
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
		},
		{
			name: "case 2 -> fail internal server error when get all task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				GetAllTaskFunc: func(ctx context.Context) ([]model.TaskModel, error) {
					return nil, errors.New("database error")
				},
			}},
			wantCode: http.StatusInternalServerError,
			wantResponse: util.ErrorResponse{
				Message: "Internal Server Error",
				Error:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCase: tt.fields.taskUseCase,
			}

			router := chi.NewRouter()
			router.Get("/api/tasks", h.GetAll)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/tasks", nil)
			router.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code, "error code")

			decoder_data := tt.wantResponse

			err := json.NewDecoder(recorder.Body).Decode(&decoder_data)
			assert.NoError(t, err, "unmarshal error")

			expect, err := json.Marshal(tt.wantResponse)
			assert.NoError(t, err, "marshal error")
			err = json.Unmarshal(expect, &tt.wantResponse)
			assert.NoError(t, err, "unmarshal error")

			assert.Equal(t, tt.wantResponse, decoder_data, "handler response")
		})
	}
}

func TestHandler_Create(t *testing.T) {
	type fields struct {
		taskUseCase *TaskUsecaseMock
	}

	type ResponseData struct {
		Message string          `json:"message"`
		Data    model.TaskModel `json:"data"`
	}

	tests := []struct {
		name         string
		fields       fields
		request      model.TaskModel
		wantCode     int
		wantResponse interface{}
	}{
		{
			name: "case 1 -> success when create task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				CreateTaskFunc: func(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
					r.ID = 1
					return r, nil
				},
			}},
			request: model.TaskModel{
				TaskName: "task 1",
				IsDone:   false,
			},
			wantCode: http.StatusCreated,
			wantResponse: ResponseData{
				Message: "Task Created",
				Data: model.TaskModel{
					ID:       1,
					TaskName: "task 1",
					IsDone:   false,
				},
			},
		},
		{
			name: "case 2 -> fail when create cause request task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				CreateTaskFunc: func(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
					return r, nil
				},
			}},
			request: model.TaskModel{
				TaskName: "",
				IsDone:   false,
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResponse: util.ErrorResponse{Message: "Invalid Request Data", Error: []model.ErrorField{
				{
					FieldName: "TaskName",
					Message:   "TaskName is required",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCase: tt.fields.taskUseCase,
			}

			router := chi.NewRouter()
			router.Post("/api/task", h.Create)
			recorder := httptest.NewRecorder()

			var buf bytes.Buffer

			err := json.NewEncoder(&buf).Encode(tt.request)
			assert.NoError(t, err, "request to buff error")

			request, _ := http.NewRequest("POST", "/api/task", &buf)
			router.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code, "error code")

			jsonExpect, err := json.Marshal(tt.wantResponse)

			assert.NoError(t, err, "unmarshal error")
			assert.Equal(t, jsonExpect, recorder.Body.Bytes(), "handler response")
		})
	}
}

func TestHandler_Update(t *testing.T) {
	type fields struct {
		taskUseCase *TaskUsecaseMock
	}

	type ResponseData struct {
		Message string          `json:"message"`
		Data    model.TaskModel `json:"data"`
	}

	tests := []struct {
		name         string
		fields       fields
		redId        interface{}
		request      model.TaskModel
		wantCode     int
		wantResponse interface{}
	}{
		{
			name: "case 1 -> success when update task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				UpdateTaskFunc: func(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
					return r, nil
				},
			}},
			redId: 1,
			request: model.TaskModel{
				TaskName: "task 1",
				IsDone:   false,
			},
			wantCode: http.StatusOK,
			wantResponse: ResponseData{
				Message: "Task Updated",
				Data: model.TaskModel{
					ID:       1,
					TaskName: "task 1",
					IsDone:   false,
				},
			},
		},
		{
			name: "case 2 -> failed when update cause request task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				UpdateTaskFunc: func(ctx context.Context, r model.TaskModel) (model.TaskModel, error) {
					return r, nil
				},
			}},
			redId: 1,
			request: model.TaskModel{
				TaskName: "",
				IsDone:   false,
			},
			wantCode: http.StatusUnprocessableEntity,
			wantResponse: util.ErrorResponse{Message: "Invalid Request Data", Error: []model.ErrorField{
				{
					FieldName: "TaskName",
					Message:   "TaskName is required",
				},
			}},
		},
		{
			name: "case 2 -> failed when update cause not found task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				UpdateTaskFunc: func(ctx context.Context, request model.TaskModel) (model.TaskModel, error) {
					return request, nil
				},
			}},
			redId: "1s",
			request: model.TaskModel{
				TaskName: "cek",
				IsDone:   false,
			},
			wantCode:     http.StatusNotFound,
			wantResponse: util.ErrorResponse{Message: "To do list not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCase: tt.fields.taskUseCase,
			}

			router := chi.NewRouter()
			router.Put("/api/task/{id}", h.Update)
			recorder := httptest.NewRecorder()

			var buf bytes.Buffer

			err := json.NewEncoder(&buf).Encode(tt.request)
			assert.NoError(t, err, "request to buff error")

			reflect.TypeOf(tt.redId)

			id := fmt.Sprintf("%v", tt.redId)

			request, _ := http.NewRequest("PUT", "/api/task/"+id, &buf)
			router.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code, "error code")

			jsonExpect, err := json.Marshal(tt.wantResponse)
			assert.NoError(t, err, "marshal error")

			assert.Equal(t, jsonExpect, recorder.Body.Bytes(), "handler response")
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type fields struct {
		taskUseCase *TaskUsecaseMock
	}

	type ResponseData struct {
		Message string        `json:"message"`
		Data    StatusRespose `json:"data"`
	}

	tests := []struct {
		name         string
		fields       fields
		redId        interface{}
		wantCode     int
		wantResponse interface{}
	}{
		{
			name: "case 1 -> success when update task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				DeleteTaskFunc: func(ctx context.Context, r model.TaskModel) error {
					return nil
				},
			}},
			redId:    1,
			wantCode: http.StatusOK,
			wantResponse: ResponseData{
				Message: "Task Deleted",
				Data: StatusRespose{
					Success: true,
				},
			},
		},
		{
			name: "case 2 -> fail when update cause already deleted task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				DeleteTaskFunc: func(ctx context.Context, r model.TaskModel) error {
					return errors.New("error, delete task failed")
				},
			}},
			redId:    1,
			wantCode: http.StatusOK,
			wantResponse: ResponseData{
				Message: "error, delete task failed",
				Data: StatusRespose{
					Success: false,
				},
			},
		},
		{
			name: "case 3 -> fail when update cause invalid id not found task handler",
			fields: fields{taskUseCase: &TaskUsecaseMock{
				DeleteTaskFunc: func(ctx context.Context, r model.TaskModel) error {
					return errors.New("error, delete task failed")
				},
			}},
			redId:    "1s",
			wantCode: http.StatusNotFound,
			wantResponse: util.ErrorResponse{
				Message: "To do list not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCase: tt.fields.taskUseCase,
			}

			router := chi.NewRouter()
			router.Delete("/api/task/{id}", h.Delete)
			recorder := httptest.NewRecorder()

			id := fmt.Sprintf("%v", tt.redId)

			request, _ := http.NewRequest("DELETE", "/api/task/"+id, nil)
			router.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code, "error code")

			jsonExpect, err := json.Marshal(tt.wantResponse)
			assert.NoError(t, err, "marshal error")
			assert.Equal(t, jsonExpect, recorder.Body.Bytes(), "handler response")
		})
	}
}
