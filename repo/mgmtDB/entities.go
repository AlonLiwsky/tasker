package mgmtDB

import "github.com/tasker/entities"

type dbStep struct {
	ID          int
	Type        string
	Params      map[string]string
	FailureStep *int
	Position    *int
}

func (s dbStep) toStep() entities.Step {
	return entities.Step{
		ID:     s.ID,
		Type:   entities.StepType(s.Type),
		Params: s.Params,
	}
}
