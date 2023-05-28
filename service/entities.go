package service

import "time"

type Task struct {
	ID    int
	Name  string
	Steps []Step
}

type StepType string

//TODO: create const

type Step struct {
	ID          int
	Task        *Task
	Type        StepType
	Params      string
	FailureStep *Step //TBD: []Step
	//A failure step should be executed by a different function that handles it owns errors and retries preventing infinite loops
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
