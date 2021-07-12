//go:generate mockgen -source ../user/password.go -destination ../user/mock/mock_password.go

package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

type updatePassword interface {
	Update(ctx context.Context, id int, new, old string) error
}

// UpdatePassword is a update user password endpoint struct
type UpdatePassword struct {
	do   updatePassword
	resp *web.Response
}

func newUpdatePassword(r *web.Response, u updatePassword) *UpdatePassword {
	return &UpdatePassword{
		do:   u,
		resp: r,
	}
}

// Do is getting user's id from URL and updates user's password
func (u *UpdatePassword) Do(r *web.Request) {
	ctx := r.Context()

	id := r.GetPathParamsInt(pathParamUserID)
	if id == nil {
		u.resp.BadRequest(ctx, fmt.Errorf("user id is mising"))
		return
	}

	var reqBody entity.PaswordRequest

	err := r.UnmarshalBodyJSON(&reqBody)
	if err != nil {
		u.resp.BadRequest(ctx, err)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		u.resp.BadRequest(ctx, err)
		return
	}

	err = u.do.Update(ctx, *id, reqBody.New, reqBody.Old)
	if errors.Is(err, entity.ErrInvalidPassword) {
		u.resp.BadRequest(ctx, err)
		return
	}
	if errors.Is(err, entity.ErrNotFound) {
		u.resp.NotFound(ctx, err)
		return
	}
	if err != nil {
		u.resp.InternalServerError(ctx, err)
		return
	}

	u.resp.Ok(ctx)
}
