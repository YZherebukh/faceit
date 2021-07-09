//go:generate mockgen -source ../country/country.go -destination ../country/mock/mock_country.go

package country

import (
	"context"

	"github.com/faceit/test/entity"
)

type client interface {
	All(ctx context.Context) ([]entity.Country, error)
}
type Country struct {
	client client
}

func New(c client) *Country {
	return &Country{
		client: c,
	}
}

func (c *Country) All(ctx context.Context) ([]entity.Country, error) {
	return c.client.All(ctx)
}
