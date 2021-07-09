package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Request is a endpoint request struct
type Request struct {
	req *http.Request
}

// NewRequest creates new request instance
func NewRequest(r *http.Request) *Request {
	return &Request{
		req: r,
	}
}

func (r *Request) Context() context.Context {
	return r.req.Context()
}

// UnmarshalBodyJSON is unmarshalling req.body into v(should be a pointer)
func (r *Request) UnmarshalBodyJSON(v interface{}) error {
	err := json.NewDecoder(r.req.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body, error: %w", err)
	}

	return nil
}

// GetPathParamsString is getting Path parameters of type string
func (r *Request) GetPathParamsString(key string) string {
	return r.req.FormValue(key)
}

// GetPathParamsInt is getting Path parameters of type int
func (r *Request) GetPathParamsInt(key string) *int {
	return getParamInt(r.req.FormValue(key))
}

// GetQueryParamsString is getting Query parameters of type string
func (r *Request) GetQueryParamsString(key string) string {
	return r.req.URL.Query().Get(key)
}

// GetQueryParamsInt is getting Query parameters of type int
func (r *Request) GetQueryParamsInt(key string) *int {
	return getParamInt(r.req.URL.Query().Get(key))
}

func getParamInt(param string) *int {
	n, err := strconv.Atoi(param)
	if err != nil {
		return nil
	}

	return &n
}
