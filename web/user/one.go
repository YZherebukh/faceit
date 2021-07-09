//go:generate mockgen -source ../user/one.go -destination ../user/mock/mock_one.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

type one interface {
	One(ctx context.Context, id int) (entity.User, error)
}

type One struct {
	do   one
	resp *web.Response
}

func newOne(r *web.Response, o one) *One {
	return &One{
		do:   o,
		resp: r,
	}
}

func (o *One) Do(r *web.Request) {
	ctx := r.Context()

	id := r.GetPathParamsInt(pathParamUserID)
	if id == nil {
		o.resp.BadRequest(ctx, entity.ErrUserIDIsMissing)
		return
	}

	user, err := o.do.One(ctx, *id)
	if errors.Is(err, entity.ErrNotFound) {
		o.resp.NotFound(ctx, err)
		return
	}

	if err != nil {
		o.resp.InternalServerError(ctx, err)
		return
	}

	o.resp.Ok(ctx).WithBody(ctx, user.ToResponse())
}
