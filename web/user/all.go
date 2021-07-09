//go:generate mockgen -source ../user/all.go -destination ../user/mock/mock_all.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

const (
	queryParamTitleConst  = "title"
	queryParamFilterConst = "filter"
)

type all interface {
	All(ctx context.Context, filter, title string) ([]entity.User, error)
}

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

func (a *All) Do(r *web.Request) {
	ctx := r.Context()

	users, err := a.do.All(ctx, r.GetPathParamsString(queryParamFilterConst), r.GetPathParamsString(queryParamTitleConst))
	if errors.Is(err, entity.ErrNotFound) {
		a.resp.NoContent(ctx)
		return
	}

	if err != nil {
		a.resp.InternalServerError(ctx, err)
		return
	}

	resp := make([]entity.UserResponse, len(users))
	for i := range users {
		resp[i] = users[i].ToResponse()
	}

	a.resp.Ok(ctx).WithBody(ctx, resp)
}
