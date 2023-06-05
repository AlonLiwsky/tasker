package entities

import (
	"errors"
	"time"

	"github.com/tasker/http"
)

type Task struct {
	ID    int
	Name  string
	Steps []Step
}

func (t Task) IsValid() error {
	if t.Name == "" {
		return http.WrapError(errors.New("task must have a name"), http.ErrBadRequest)
	}

	if len(t.Steps) == 0 {
		return http.WrapError(errors.New("task must have steps"), http.ErrBadRequest)
	}

	for _, step := range t.Steps {
		if err := step.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

type StepType string

const (
	APICallStepType StepType = "api_call"
)

type Step struct {
	ID          int
	Type        StepType
	Params      map[string]string
	FailureStep *Step `json:"failure_step"`
	//A failure step should be executed by a different function that handles it owns errors and retries preventing infinite loops
}

func (s Step) IsValid() error {
	validTypes := map[StepType]bool{
		APICallStepType: true,
	}
	if !validTypes[s.Type] {
		return http.WrapError(errors.New("step must have a valid step type"), http.ErrBadRequest)
	}

	if len(s.Params) == 0 {
		return http.WrapError(errors.New("step must have a params"), http.ErrBadRequest)
	}

	return nil
}

type ScheduledTask struct {
	ID          int
	Name        string
	Chron       string
	RetryPolicy string
	Task        *Task
	Enabled     bool
	LastRun     time.Time
	FirstRun    time.Time
}

type Execution struct {
	ID                   int
	ScheduledTask        *ScheduledTask
	TryNumber            int
	Status               string
	RequestedTime        time.Time
	ExecutedTime         time.Time
	LastStatusChangeTime time.Time
}