//go:generate mockgen -source ../user/user.go -destination ../user/mock/mock_user.go

package user

import (
	"context"
	"errors"

	"github.com/faceit/test/entity"
)

const (
	filterCountry   = "country"
	filterLastName  = "first_name"
	filterFirstName = "last_name"
	filterNickname  = "nick_name"
	filterEmil      = "email"
)

type client interface {
	Create(ctx context.Context, u entity.User) (int, error)
	Update(ctx context.Context, u entity.User) error
	Delete(ctx context.Context, id int) error
	One(ctx context.Context, id int) (entity.User, error)
	All(ctx context.Context) ([]entity.User, error)
	AllByCountry(ctx context.Context, iso2 string) ([]entity.User, error)
	AllWithFilter(ctx context.Context, title, filter string) ([]entity.User, error)
}

type hasher interface {
	HashAndSalt(password string) error
	Salt() string
	Hashed() string
}

type User struct {
	client client
	hasher hasher
}

func New(c client, h hasher) *User {
	return &User{
		client: c,
		hasher: h,
	}
}

func (u *User) Create(ctx context.Context, user entity.User) (int, error) {
	_, err := u.One(ctx, user.ID)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return 0, err
	}

	err = u.hasher.HashAndSalt(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = u.hasher.Hashed()
	user.Salt = u.hasher.Salt()

	return u.client.Create(ctx, user)
}

func (u *User) All(ctx context.Context, filter, title string) ([]entity.User, error) {
	switch filter {
	case filterCountry:
		return u.client.AllByCountry(ctx, title)
	case filterFirstName:
		fallthrough
	case filterLastName:
		fallthrough
	case filterNickname:
		fallthrough
	case filterEmil:
		return u.client.AllWithFilter(ctx, filter, title)
	default:
		return u.client.All(ctx)
	}
}

func (u *User) One(ctx context.Context, id int) (entity.User, error) {
	return u.client.One(ctx, id)
}

func (u *User) Update(ctx context.Context, user entity.User) error {
	err := u.CanUpdate(ctx, user)
	if err != nil {
		return err
	}

	return u.client.Update(ctx, user)
}

func (u *User) Delete(ctx context.Context, user entity.User) error {
	err := u.CanUpdate(ctx, user)
	if err != nil {
		return err
	}

	return u.client.Delete(ctx, user.ID)
}

func (u *User) CanUpdate(ctx context.Context, user entity.User) error {
	original, err := u.client.One(ctx, user.ID)
	if err != nil {
		return err
	}

	err = u.hasher.HashAndSalt(user.Password)
	if err != nil {
		return err
	}

	if original.Password != u.hasher.Hashed() {
		return entity.ErrValidationFailed
	}

	return nil
}
