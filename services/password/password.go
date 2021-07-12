//go:generate mockgen -source ../password/password.go -destination ../password/mock/mock_password.go

package password

import (
	"context"
	"fmt"

	"github.com/faceit/test/entity"
)

// client is a password client interface
type client interface {
	Update(ctx context.Context, userID int, hash, salt string) error
	One(ctx context.Context, id int) (entity.Password, error)
}

// hasher is a password hasher interface
type hasher interface {
	Hash(password, salt string) (string, error)
	Salt() string
	Compare(password, hashed string) error
}

// Password is a password service struct
type Password struct {
	client
	hasher
}

// New creates new password service
func New(c client, h hasher) *Password {
	return &Password{
		client: c,
		hasher: h,
	}
}

// Update updates user password
func (p *Password) Update(ctx context.Context, id int, new, old string) error {
	pass, err := p.client.One(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user's password from store, error: %w", err)
	}

	err = p.hasher.Compare(old, pass.Hash)
	if err != nil {
		return err
	}

	salt := p.hasher.Salt()
	newHashed, err := p.hasher.Hash(new, salt)
	if err != nil {
		return err
	}

	err = p.client.Update(ctx, id, newHashed, salt)
	if err != nil {
		return fmt.Errorf("failed to update user's password, error: %w", err)
	}

	return nil
}
