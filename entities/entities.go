package entities

import (
	"errors"
	"time"

	"github.com/tasker/http"
)

type Task struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
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

func GetAllStepTypes() []StepType {
	return []StepType{
		APICallStepType,
	}
}

type Step struct {
	ID          int               `json:"id"`
	Type        StepType          `json:"type"`
	Params      map[string]string `json:"params"`
	FailureStep *Step             `json:"failure_step"`
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

	//check for nested failure steps
	if s.FailureStep != nil && s.FailureStep.FailureStep != nil {
		return http.WrapError(errors.New("a failure step cant have its own failure step"), http.ErrBadRequest)
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

type executionStatus string

const (
	SuccessExecutionStatus        = executionStatus("success")
	FailureExecutionStatus        = executionStatus("failure")
	HandledFailureExecutionStatus = executionStatus("handled_failure")
)

type Execution struct {
	ID            int `json:"id"`
	ScheduledTask int `json:"scheduled_task"`
	//TryNumber            int
	Status executionStatus `json:"status"`
	//RequestedTime        time.Time
	ExecutedTime time.Time `json:"executed_time"`
	//LastStatusChangeTime time.Time
}

/*
Who makes the retries?
if it's the one calling execute, execution doesn't need TryNumber, RequestedTime and LastStatusChangeTime
*/
