package service

import (
	"errors"
	"time"
)

type Task struct {
	ID    int
	Name  string
	Steps []Step
}

func (t Task) IsValid() error {
	if t.Name == "" {
		return errors.New("task must have a name")
	}

	if len(t.Steps) == 0 {
		return errors.New("task must have steps")
	}

	for _, step := range t.Steps {
		if err := step.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

type StepType string

//TODO: create const

type Step struct {
	ID int
	//Task        *Task
	Type        StepType
	Params      map[string]string
	FailureStep *Step `json:"failure_step"`
	//A failure step should be executed by a different function that handles it owns errors and retries preventing infinite loops
}

func (s Step) IsValid() error {
	if s.Type == "" {
		return errors.New("step must have a type")
	}

	if len(s.Params) == 0 {
		return errors.New("step must have a params")
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
