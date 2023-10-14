package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	model "to-do-list/internal/model/task"
	util "to-do-list/pkg/response"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	useCase TaskUsecase
}

type ResponseStandard struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StatusRespose struct {
	Success bool `json:"status"`
}

func NewHandler(useCase TaskUsecase) *Handler {
	return &Handler{useCase: useCase}
}

type TaskUsecase interface {
	GetAllTask(ctx context.Context) ([]model.TaskModel, error)
	CreateTask(ctx context.Context, r model.TaskModel) (model.TaskModel, error)
	UpdateTask(ctx context.Context, r model.TaskModel) (model.TaskModel, error)
	DeleteTask(ctx context.Context, r model.TaskModel) error
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responses, err := h.useCase.GetAllTask(ctx)

	if err != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{Message: "Internal Server Error"}, http.StatusInternalServerError, w)
	}

	if err := util.ResponseJSON(responses, http.StatusOK, w); err != nil {
		fmt.Println("[Get All] Response error")
	}

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request = model.TaskModel{}
		status  = http.StatusCreated
	)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{Message: "Invalid Request Data"}, http.StatusUnprocessableEntity, w)
		return
	}
	json.Unmarshal(reqBody, &request)

	validate := Validate(request)
	if validate != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{
			Message: "Invalid Request Data",
			Error:   validate,
		}, http.StatusUnprocessableEntity, w)
		return
	}

	data, err := h.useCase.CreateTask(ctx, request)

	responses := ResponseStandard{
		Message: "Task Created",
		Data:    data,
	}

	if err != nil {
		status = http.StatusInternalServerError
		responses.Message = err.Error()
		responses.Data = StatusRespose{Success: false}
	}

	if err := util.ResponseJSON(responses, status, w); err != nil {
		fmt.Println("[Create] Response error")
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request = model.TaskModel{}
	)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{Message: "Invalid Request Data"}, http.StatusUnprocessableEntity, w)
		return
	}
	json.Unmarshal(reqBody, &request)

	validate := Validate(request)
	if validate != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{
			Message: "Invalid Request Data",
			Error:   validate,
		}, http.StatusUnprocessableEntity, w)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{Message: "To do list not found"}, http.StatusNotFound, w)
		return
	}

	request.ID = id
	data, err := h.useCase.UpdateTask(ctx, request)

	responses := ResponseStandard{
		Message: "Task Updated",
		Data:    data,
	}

	if err != nil {
		responses.Message = err.Error()
		responses.Data = StatusRespose{Success: false}
	}

	if err := util.ResponseJSON(responses, http.StatusOK, w); err != nil {
		fmt.Println("[Update] Response error")
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		util.ResponseErrorJSON(&util.ErrorResponse{Message: "To do list not found"}, http.StatusNotFound, w)
		return
	}

	status_response := StatusRespose{Success: true}

	responses := ResponseStandard{
		Message: "Task Deleted",
		Data:    status_response,
	}

	request := model.TaskModel{ID: id}
	err = h.useCase.DeleteTask(ctx, request)
	if err != nil {
		responses.Message = err.Error()
		responses.Data = StatusRespose{Success: false}
	}

	if err := util.ResponseJSON(responses, http.StatusOK, w); err != nil {
		fmt.Println("[Update] Response error")
	}
}
