package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tasker/service"
)

type Service interface {
	CreateTask(ctx context.Context, task service.Task) (service.Task, error)
}

type adapter struct {
	service Service
}

func (a adapter) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedTask := service.Task{}
	if err := decode(r, &receivedTask); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request. There is an error in your JSON body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	task, err := a.service.CreateTask(ctx, receivedTask)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding task as JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf(`{"msg": "task saved successfully", "task": %s}`, taskJSON)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func NewAdapter(srv Service) *adapter {
	return &adapter{service: srv}
}
