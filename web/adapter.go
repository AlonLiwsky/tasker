package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/tasker/entities"
	httpErr "github.com/tasker/http"
)

type Service interface {
	CreateTask(ctx context.Context, task entities.Task) (entities.Task, error)
	GetTask(ctx context.Context, taskID int) (entities.Task, error)
	ExecuteTask(ctx context.Context, taskID int, scheduleID int, idempToken string) (entities.Execution, error)
	CreateSchedule(ctx context.Context, sch entities.ScheduledTask) (entities.ScheduledTask, error)
	ExecuteScheduledTasks(ctx context.Context) error
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

func (a adapter) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest.WithMessage("invalid task ID")))
		return
	}

	task, err := a.service.GetTask(ctx, taskID)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskJSON)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}
}

func (a adapter) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest.WithMessage("invalid task ID")))
		return
	}

	schID, err := strconv.Atoi(chi.URLParam(r, "scheduleID"))
	if err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest.WithMessage("invalid schedule ID")))
		return
	}

	idempotencyTokenMsg := struct {
		Token string `json:"idempotency_token"`
	}{}
	if err := decode(r, &idempotencyTokenMsg); err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest.WithMessage("invalid idempotency token")))
		return
	}
	if idempotencyTokenMsg.Token == "" {
		httpErr.JSONHandleError(w, httpErr.ErrBadRequest.WithMessage("invalid idempotency token"))
		return
	}

	execution, err := a.service.ExecuteTask(ctx, taskID, schID, idempotencyTokenMsg.Token)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	execJSON, err := json.Marshal(execution)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(execJSON)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}
}

func (a adapter) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedSchedule := ScheduledTask{}
	if err := decode(r, &receivedSchedule); err != nil {
		httpErr.JSONHandleError(w, httpErr.WrapError(err, httpErr.ErrBadRequest))
		return
	}
	sch := entities.ScheduledTask{
		Name:    receivedSchedule.Name,
		Cron:    receivedSchedule.Cron,
		Retries: receivedSchedule.Retries,
		Task:    entities.Task{ID: receivedSchedule.TaskID},
		Enabled: receivedSchedule.Enabled,
	}

	if err := sch.IsValid(); err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	sch, err := a.service.CreateSchedule(ctx, sch)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	schJSON, err := json.Marshal(sch)
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf(`{"msg": "schedule saved successfully", "schedule": %s}`, schJSON)))
	if err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}
}

func (a adapter) ExecuteScheduledTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := a.service.ExecuteScheduledTasks(ctx); err != nil {
		httpErr.JSONHandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("scheduled tasks execution finished successfully"))
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
