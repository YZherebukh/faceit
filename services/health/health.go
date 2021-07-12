package health

import (
	"context"
	"time"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/logger"
)

type check interface {
	Ping() error
}

// Check is a healtCheck struct
type Check struct {
	db  check
	log logger.Logger
}

// New creates New Check instance
func New(db check, log logger.Logger) *Check {
	return &Check{
		db: db,
	}
}

// Do performes a healtCheck
func (c *Check) Do(ctx context.Context) []entity.Response {
	return []entity.Response{
		c.dbCheck(ctx),
	}
}

// dbCheck is checking db health
func (c *Check) dbCheck(ctx context.Context) entity.Response {
	resp := entity.Response{
		Name: "database",
		Time: time.Now().UTC().Format(time.RFC3339),
	}

	err := c.db.Ping()
	if err != nil {
		resp.Message = err.Error()

		c.log.Errorf(ctx, "healthCheck: %s in unhealthy, error: %w", err)
	}

	resp.Healthy = err == nil

	return resp
}
