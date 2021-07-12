//go:generate mockgen -source ../user/all.go -destination ../user/mock/mock_all.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

// query params
const (
	queryParamTitleConst  = "title"
	queryParamFilterConst = "filter"
)

type all interface {
	All(ctx context.Context, title, filter string) ([]entity.User, error)
}

// All is a all users endpoint struct
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

// Do returnes all users based on provided filters,
// if o filters provided, all users would be returned
func (a *All) Do(r *web.Request) {
	ctx := r.Context()

	users, err := a.do.All(ctx, r.GetQueryParamsString(queryParamTitleConst), r.GetQueryParamsString(queryParamFilterConst))
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
