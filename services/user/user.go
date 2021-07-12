//go:generate mockgen -source ../user/user.go -destination ../user/mock/mock_user.go

package user

import (
	"context"

	"github.com/faceit/test/entity"
)

// user get all filter title
const (
	filterCountry   = "country"
	filterLastName  = "firstName"
	filterFirstName = "lastName"
	filterNickname  = "nickName"
	filterEmil      = "email"
)

var (
	titleMap = map[string]string{
		filterLastName:  "first_name",
		filterFirstName: "last_name",
		filterNickname:  "nick_name",
		filterEmil:      "email",
	}
)

// client is a user client interface
type client interface {
	Create(ctx context.Context, u entity.User) (int, error)
	Update(ctx context.Context, u entity.User) error
	One(ctx context.Context, id int) (entity.User, error)
	Delete(ctx context.Context, id int) error
	All(ctx context.Context) ([]entity.User, error)
	AllByCountry(ctx context.Context, iso2 string) ([]entity.User, error)
	AllWithFilter(ctx context.Context, title, filter string) ([]entity.User, error)
}

type passwordClient interface {
	One(ctx context.Context, id int) (entity.Password, error)
}

// hasher is a user password hasher interface
type hasher interface {
	Hash(password, salt string) (string, error)
	Salt() string
	Compare(password, hashed string) error
}

// user is a user service struct
type User struct {
	client         client
	passwordClient passwordClient
	hasher         hasher
}

// New creates new user service instance
func New(c client, h hasher, p passwordClient) *User {
	return &User{
		client:         c,
		passwordClient: p,
		hasher:         h,
	}
}

// Create creates new user in store
func (u *User) Create(ctx context.Context, user entity.User) (int, error) {
	user.Salt = u.hasher.Salt()

	pass, err := u.hasher.Hash(user.Password, user.Salt)
	if err != nil {
		return 0, err
	}

	user.Password = pass

	return u.client.Create(ctx, user)
}

// All returnes all users from store depending on filter and title
func (u *User) All(ctx context.Context, title, filter string) ([]entity.User, error) {
	switch title {
	case filterCountry:
		return u.client.AllByCountry(ctx, filter)
	case filterFirstName:
		fallthrough
	case filterLastName:
		fallthrough
	case filterNickname:
		fallthrough
	case filterEmil:
		return u.client.AllWithFilter(ctx, titleMap[title], filter)
	default:
		return u.client.All(ctx)
	}
}

// One returnes one user b ID
func (u *User) One(ctx context.Context, id int) (entity.User, error) {
	return u.client.One(ctx, id)
}

// Update updates user by ID
func (u *User) Update(ctx context.Context, user entity.User) error {
	err := u.canUpdate(ctx, user)
	if err != nil {
		return err
	}

	return u.client.Update(ctx, user)
}

// Delete deletes user by id
func (u *User) Delete(ctx context.Context, user entity.User) error {
	err := u.canUpdate(ctx, user)
	if err != nil {
		return err
	}

	return u.client.Delete(ctx, user.ID)
}

// Can update is checking if action on user can be perfrmed by comparing
// existing password with provided
func (u *User) canUpdate(ctx context.Context, user entity.User) error {
	original, err := u.passwordClient.One(ctx, user.ID)
	if err != nil {
		return err
	}

	return u.hasher.Compare(user.Password, original.Hash)
}
