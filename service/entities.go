package service

import "time"

type Task struct {
	ID        int
	Name      string
	FirstStep *Step
}

type Step struct {
	ID              int
	Task            *Task
	SuccessNextStep *Step
	FailureNextStep *Step
	StepID          int
	Params          string
	FinalStep       bool
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
