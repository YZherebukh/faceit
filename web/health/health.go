package health

import (
	"context"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

// health is a Health endpoint interface
type health interface {
	Do(ctx context.Context) []entity.Response
}

// Health is a health endpoint struct
type Health struct {
	health health
	resp   *web.Response
}

// newHealth creates new Health instance
func newHealth(r *web.Response, h health) *Health {
	return &Health{
		health: h,
		resp:   r,
	}
}

// Do is performing a health check
func (h *Health) Do(r *web.Request) {
	ctx := r.Context()

	resp := h.health.Do(ctx)

	h.resp.WithBody(ctx, resp).Ok(ctx)
}
