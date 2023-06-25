package apicall

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tasker/service/apicall"
)

func (r repository) ApiCall(ctx context.Context, method, url, body string, headers map[string][]string) (apicall.Response, error) {
	// Parse body as
	httpBody := strings.NewReader(body)

	// Make request
	request, err := http.NewRequestWithContext(ctx, method, url, httpBody)
	if err != nil {
		return apicall.Response{}, fmt.Errorf("preparing API request to %s: %w", url, err)
	}
	resp, err := r.client.Do(request)
	if err != nil {
		return apicall.Response{}, fmt.Errorf("making API call to %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read the Response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return apicall.Response{}, fmt.Errorf("failed to read response body from API call to %s: %w", url, err)
	}

	return apicall.Response{
		Body:       string(responseBody),
		StatusCode: resp.StatusCode,
	}, nil
}
