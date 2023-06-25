package storageread

import (
	"context"
	"fmt"
)

const (
	storageKeyParam = "storage_key"
)

type Repository interface {
	Get(ctx context.Context, key string) (string, error)
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
		return "", fmt.Errorf("no key param found for storage read step")
	}

	return a.repo.Get(ctx, key)
}
