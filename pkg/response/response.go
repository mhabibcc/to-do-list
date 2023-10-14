package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string
	Error   interface{}
}

func ResponseJSON(data interface{}, status int, writer http.ResponseWriter) (err error) {
	writer.Header().Set("Content-type", "aplication/json")

	writer.WriteHeader(status)

	d, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		d, _ = json.Marshal(ErrorResponse{Message: "ResponseJSON: Failed to response " + err.Error()})
		err = fmt.Errorf("ResponseJSON: Failed to response : %s", err)
	}

	_, _ = writer.Write(d)

	return
}

func ResponseErrorJSON(errResponse *ErrorResponse, status int, writer http.ResponseWriter) (err error) {
	writer.Header().Set("Content-type", "aplication/json")

	writer.WriteHeader(status)
	d, _ := json.Marshal(errResponse)
	_, _ = writer.Write(d)

	return
}
