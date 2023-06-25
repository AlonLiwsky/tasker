package storagewrite

import (
	"context"
	"fmt"

	"github.com/tasker/service"
)

const (
	storageKeyParam   = "storage_key"
	storageValueParam = "storage_value"
)

type Repository interface {
	Set(ctx context.Context, key, value string) error
}

type stepRunner struct {
	repo Repository
}

func NewStepRunner(repo Repository) stepRunner {
	return stepRunner{repo: repo}
}

func (a stepRunner) RunStep(ctx context.Context, params map[string]string) (string, error) {
	// Get key from params
	key, found := params[storageKeyParam]
	if !found {
		return "", fmt.Errorf("no key param found for storage write step")
	}
	if key == service.UseLastStepResultKey {
		key, found = params[service.LastStepResultKey]
		if !found {
			return "", fmt.Errorf("requested to use last step result as key but there was no last step result")
		}
	}

	// Get value from params
	value, found := params[storageValueParam]
	if !found {
		return "", fmt.Errorf("no value param found for storage write step")
	}
	if value == service.UseLastStepResultKey {
		value, found = params[service.LastStepResultKey]
		if !found {
			return "", fmt.Errorf("requested to use last step result as value but there was no last step result")
		}
	}

	return value, a.repo.Set(ctx, key, value)
}
