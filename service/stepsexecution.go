package service

import (
	"context"
	"fmt"

	"github.com/tasker/entities"
)

type StepRunner interface {
	RunStep(ctx context.Context, params map[string]string) (string, error)
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

func (s service) runStep(ctx context.Context, step entities.Step) (string, error) {
	return s.stepRunners[step.Type].RunStep(ctx, step.Params)
}
