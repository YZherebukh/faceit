//go:generate mockgen -source ../user/update.go -destination ../user/mock/mock_update.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

const (
	pathParamUserID = "id"
	actionUpdate    = "UPDATE"
)

type update interface {
	Update(ctx context.Context, u entity.User) error
}

type Update struct {
	do        update
	resp      *web.Response
	notify    notifier
	consumers []string
}

func newUpdate(r *web.Response, u update, n notifier, consumers []string) *Update {
	return &Update{
		do:        u,
		resp:      r,
		notify:    n,
		consumers: consumers,
	}
}

func (u *Update) Do(r *web.Request) {
	ctx := r.Context()

	id := r.GetPathParamsInt(pathParamUserID)
	if id == nil {
		u.resp.BadRequest(ctx, entity.ErrUserIDIsMissing)
		return
	}

	var reqBody entity.UserRequest

	err := r.UnmarshalBodyJSON(&reqBody)
	if err != nil {
		u.resp.BadRequest(ctx, err)
		return
	}

	user := reqBody.ToUser()
	user.ID = *id

	err = u.do.Update(ctx, user)
	if errors.Is(err, entity.ErrNotFound) {
		u.resp.NotFound(ctx, err)
		return
	}
	if err != nil {
		u.resp.InternalServerError(ctx, err)
		return
	}

	u.notify.Add(entity.NotifierMessage{
		Message: entity.UserNotification{
			User:   user.ToResponse(),
			Action: actionUpdate},
		Consumers: u.consumers})

	u.resp.Ok(ctx)
}
