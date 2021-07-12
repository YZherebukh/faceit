package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/faceit/test/logger"
)

// header constants
const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=UTF-8"
)

// Response is an endpoint response struct
type Response struct {
	writer http.ResponseWriter
	log    logger.Logger
}

// NewResponse creates new Response instance
func NewResponse(w http.ResponseWriter, log logger.Logger) *Response {
	return &Response{
		log:    log,
		writer: w,
	}
}

// ContentHeader is setting a Content-Type header to "application/json; charset=UTF-8"
func (r *Response) ContentHeader(ctx context.Context) *Response {
	r.log.Infof(ctx, "setting headers %s:%s", contentTypeKey, contentTypeValue)

	r.writer.Header().Set(contentTypeKey, contentTypeValue)

	return r
}

// Ok is setting response status code to http.StatusOK
func (r *Response) Ok(ctx context.Context) *Response {
	return r.setStatus(ctx, http.StatusOK)
}

// Created is setting response status code to http.StatusCreated
func (r *Response) Created(ctx context.Context) *Response {
	return r.setStatus(ctx, http.StatusCreated)
}

// NoContent is setting response status code to http.StatusNoContent
func (r *Response) NoContent(ctx context.Context) *Response {
	return r.setStatus(ctx, http.StatusNoContent)
}

// Unauthorized is setting response status code to http.StatusUnauthorized
func (r *Response) Unauthorized(ctx context.Context) *Response {
	return r.setStatus(ctx, http.StatusUnauthorized)
}

// NotFound is setting response status code to http.StatusUnauthorized
func (r *Response) NotFound(ctx context.Context, err error) *Response {
	r.log.Warningf(ctx, "bad request, message: %s", err.Error())

	return r.setStatus(ctx, http.StatusNotFound)
}

// BadRequest is setting response status code to http.StatusBadRequest
func (r *Response) BadRequest(ctx context.Context, err error) *Response {
	r.log.Warningf(ctx, "bad request, message: %s", err.Error())

	return r.setStatus(ctx, http.StatusBadRequest)
}

// InternalServerError is setting response status code to http.InternalServerError
func (r *Response) InternalServerError(ctx context.Context, err error) *Response {
	r.log.Errorf(ctx, "request failed, error: %s", err.Error())

	return r.setStatus(ctx, http.StatusInternalServerError)
}

// Ok is marshalling v and sets it as a response body
func (r *Response) WithBody(ctx context.Context, v interface{}) *Response {
	body, err := json.Marshal(v)
	if err != nil {
		r.log.Errorf(ctx, "failed to marsahl response body, error: %s", err.Error())
	}

	_, err = r.writer.Write(body)
	if err != nil {
		r.log.Errorf(ctx, "failed to write response body, error: %s", err.Error())
	}

	return r
}

func (r *Response) setStatus(ctx context.Context, status int) *Response {
	r.log.Infof(ctx, "sending status %d", status)

	r.writer.WriteHeader(status)

	return r
}
