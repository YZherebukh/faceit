//go:generate mockgen -source ../user/create.go -destination ../user/mock/mock_create.go

package user

import (
	"context"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/web"
)

const (
	actionCreate = "CREATE"
)

type create interface {
	Create(ctx context.Context, u entity.User) (int, error)
}

type notifier interface {
	Add(message entity.NotifierMessage)
}

type Create struct {
	do        create
	resp      *web.Response
	notify    notifier
	consumers []string
}

func newCreate(r *web.Response, c create, n notifier, consumers []string) *Create {
	return &Create{
		do:        c,
		resp:      r,
		notify:    n,
		consumers: consumers,
	}
}

func (c *Create) Do(r *web.Request) {
	ctx := r.Context()

	var reqBody entity.UserRequest

	err := r.UnmarshalBodyJSON(&reqBody)
	if err != nil {
		c.resp.BadRequest(ctx, err)
		return
	}

	err = reqBody.Validate()
	if err != nil {
		c.resp.BadRequest(ctx, err)
		return
	}

	user := reqBody.ToUser()

	id, err := c.do.Create(ctx, user)
	if err != nil {
		c.resp.InternalServerError(ctx, err)
		return
	}

	user.ID = id

	c.notify.Add(entity.NotifierMessage{
		Message: entity.UserNotification{
			User:   user.ToResponse(),
			Action: actionCreate},
		Consumers: c.consumers})

	var respbody struct {
		ID int `json:"id"`
	}

	respbody.ID = id

	c.resp.Created(ctx).WithBody(ctx, respbody)
}
