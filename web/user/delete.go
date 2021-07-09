//go:generate mockgen -source ../user/delete.go -destination ../user/mock/mock_delete.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

const (
	actionDelete = "DELETE"
)

type delete interface {
	Delete(ctx context.Context, user entity.User) error
}

type Delete struct {
	do        delete
	resp      *web.Response
	notify    notifier
	consumers []string
}

func newDelete(r *web.Response, d delete, n notifier, consumers []string) *Delete {
	return &Delete{
		do:        d,
		resp:      r,
		notify:    n,
		consumers: consumers,
	}
}

func (d *Delete) Do(r *web.Request) {
	ctx := r.Context()

	id := r.GetPathParamsInt(pathParamUserID)
	if id == nil {
		d.resp.BadRequest(ctx, entity.ErrUserIDIsMissing)
		return
	}

	var reqBody entity.UserRequest

	err := r.UnmarshalBodyJSON(&reqBody)
	if err != nil {
		d.resp.BadRequest(ctx, err)
		return
	}

	deleteUser := reqBody.ToUser()
	deleteUser.ID = *id

	err = d.do.Delete(ctx, deleteUser)
	if errors.Is(err, entity.ErrNotFound) {
		d.resp.NotFound(ctx, err)
		return
	}
	if err != nil {
		d.resp.InternalServerError(ctx, err)
		return
	}

	d.notify.Add(entity.NotifierMessage{
		Message: entity.UserNotification{
			User:   deleteUser.ToResponse(),
			Action: actionDelete},
		Consumers: d.consumers})

	d.resp.Ok(ctx)
}
