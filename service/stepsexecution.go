package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

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

const (
	urlAPIStepParam         = "url_api"
	requestVerbAPIStepParam = "request_verb_api"
	headersAPIStepParam     = "headers_api"
	bodyAPIStepParam        = "body_api"
)

type apiCallerStepRunner struct {
	client http.Client
}

func NewApiCallerStepRunner(client http.Client) apiCallerStepRunner {
	return apiCallerStepRunner{client: client}
}

func (a apiCallerStepRunner) RunStep(ctx context.Context, params map[string]string) (string, error) {
	// Get URL from params
	url, found := params[urlAPIStepParam]
	if !found {
		return "", fmt.Errorf("no url param found for api call step")
	}

	// Get request verb from params
	requestVerb, found := params[requestVerbAPIStepParam]
	if !found {
		return "", fmt.Errorf("no request verb param found for api call step")
	}

	// Check validity of request verb
	switch requestVerb {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
	default:
		return "", fmt.Errorf("invalid request verb param found for api call step")
	}

	// Get request body from params
	var body io.Reader
	bodyStr, found := params[bodyAPIStepParam]
	if found {
		body = strings.NewReader(bodyStr)
	}

	// Make request
	request, err := http.NewRequestWithContext(ctx, requestVerb, url, body)
	if err != nil {
		return "", fmt.Errorf("preparing API request to %s: %w", url, err)
	}
	resp, err := a.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("making API call to %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from API call to %s: %w", url, err)
	}

	//Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return string(responseBody), fmt.Errorf("making API call to %s failed with code %d: %s", url, resp.StatusCode, responseBody)
	}

	return string(responseBody), nil
}

func (s service) runStep(ctx context.Context, step entities.Step) (string, error) {
	return s.stepRunners[step.Type].RunStep(ctx, step.Params)
}
