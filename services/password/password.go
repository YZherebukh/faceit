//go:generate mockgen -source ../password/password.go -destination ../password/mock/mock_password.go

package password

import (
	"context"

	"github.com/faceit/test/entity"
)

type client interface {
	Update(ctx context.Context, userID int, hash, salt string) error
	One(ctx context.Context, id int) (entity.Password, error)
}

type hasher interface {
	HashAndSalt(password string) error
	Salt() string
	Hashed() string
}

type Password struct {
	client
	hasher
}

func New(c client, h hasher) *Password {
	return &Password{
		client: c,
		hasher: h,
	}
}

func (p *Password) Update(ctx context.Context, id int, new, old string) error {
	err := p.hasher.HashAndSalt(new)
	if err != nil {
		return err
	}

	pass, err := p.client.One(ctx, id)
	if err != nil {
		return err
	}

	if old == pass.Hash {
		return p.client.Update(ctx, id, p.hasher.Hashed(), p.hasher.Salt())
	}

	return entity.ErrValidationFailed
}
