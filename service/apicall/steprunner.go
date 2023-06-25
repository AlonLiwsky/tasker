package apicall

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	urlParam         = "url_api"
	requestVerbParam = "request_verb_api"
	headersParam     = "headers_api"
	bodyParam        = "body_api"
)

type Response struct {
	Body       string
	StatusCode int
}

type Repository interface {
	ApiCall(ctx context.Context, method, url, body string, headers map[string][]string) (Response, error)
}

type stepRunner struct {
	repo Repository
}

func NewStepRunner(repo Repository) stepRunner {
	return stepRunner{repo: repo}
}

func (a stepRunner) RunStep(ctx context.Context, params map[string]string) (string, error) {
	// Get URL from params
	url, found := params[urlParam]
	if !found {
		return "", fmt.Errorf("no url param found for api call step")
	}

	// Get request verb from params
	requestVerb, found := params[requestVerbParam]
	if !found {
		return "", fmt.Errorf("no request verb param found for api call step")
	}

	// Check validity of request verb
	switch requestVerb {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
	default:
		return "", fmt.Errorf("invalid request verb param found for api call step")
	}

	var headers map[string][]string
	headersJSON, found := params[headersParam]
	if found {
		if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
			log.Printf("Error unmarshalling headers: %s. Executing api call without provided headers", err)
		}
	}

	resp, err := a.repo.ApiCall(ctx, requestVerb, url, params[bodyParam], headers)
	if err != nil {
		return "", fmt.Errorf("making api call: %w", err)
	}

	//Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Body, fmt.Errorf("making API call to %s failed with code %d: %s", url, resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
}
