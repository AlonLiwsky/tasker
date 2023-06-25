package apicall

import (
	"context"
	"net/http"

	"github.com/tasker/service/apicall"
)

type Repository interface {
	ApiCall(ctx context.Context, method, url, body string, headers map[string][]string) (apicall.Response, error)
}

type repository struct {
	client http.Client
}

func NewRepository(client http.Client) Repository {
	return &repository{client: client}
}
