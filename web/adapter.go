package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tasker/entities"
	httpErr "github.com/tasker/http"
)

type Service interface {
	CreateTask(ctx context.Context, task entities.Task) (entities.Task, error)
}

type adapter struct {
	service Service
}

func (a adapter) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedTask := entities.Task{}
	if err := decode(r, &receivedTask); err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest))
		return
	}

	if err := receivedTask.IsValid(); err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	task, err := a.service.CreateTask(ctx, receivedTask)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf(`{"msg": "task saved successfully", "task": %s}`, taskJSON)))
	if err != nil {
		httpErr.JSONHandleError(w, err)
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
