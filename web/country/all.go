//go:generate mockgen -source ../country/all.go -destination ../country/mock/mock_all.go

package country

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

type all interface {
	All(ctx context.Context) ([]entity.Country, error)
}

// All is a All countries endpoint struct
type All struct {
	do   all
	resp *web.Response
}

func newAll(r *web.Response, a all) *All {
	return &All{
		do:   a,
		resp: r,
	}
}

// Do returnes a list of all countries
func (a *All) Do(r *web.Request) {
	ctx := r.Context()

	countries, err := a.do.All(ctx)
	if errors.Is(err, entity.ErrNotFound) {
		a.resp.NoContent(ctx)
		return
	}
	if err != nil {
		a.resp.InternalServerError(ctx, err)
		return
	}

	a.resp.WithBody(ctx, countries).Ok(ctx)
}
