package service

import (
	"context"
	"fmt"

	"github.com/tasker/entities"
)

type StepRunner interface {
	RunStep(params map[string]string) (string, error)
}

func validStepRunners(runners map[entities.StepType]StepRunner) error {
	stepTypes := entities.GetAllStepTypes()
	for _, stepType := range stepTypes {
		if _, found := runners[stepType]; !found {
			return fmt.Errorf("%s StepType was not found on the StepRunners map", stepType)
		}
	}
	return nil
}

type ApiCallerStepRunner struct {
	//http server?
}

func (a ApiCallerStepRunner) RunStep(params map[string]string) (string, error) {
	//TODO implement me
	fmt.Printf("making request to %s", params["url"])
	return params["url"], nil
}

func (s service) runStep(ctx context.Context, step entities.Step) (string, error) {
	return s.stepRunners[step.Type].RunStep(step.Params)
}
