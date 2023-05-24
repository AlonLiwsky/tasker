package web

import (
	"encoding/json"
	"net/http"
)

type Service interface {
}

type adapter struct {
	service Service
}

func (a adapter) CreateTask(w http.ResponseWriter, r *http.Request) {

}

func decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func NewAdapter(srv Service) *adapter {
	return &adapter{service: srv}
}
