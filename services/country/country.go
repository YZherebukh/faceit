//go:generate mockgen -source ../country/country.go -destination ../country/mock/mock_country.go

package country

import (
	"context"

	"github.com/faceit/test/entity"
)

// client is a country client interface
type client interface {
	All(ctx context.Context) ([]entity.Country, error)
}

// Country is a country service struct
type Country struct {
	client client
}

// New creates New country service
func New(c client) *Country {
	return &Country{
		client: c,
	}
}

// All returnes all countries from store
func (c *Country) All(ctx context.Context) ([]entity.Country, error) {
	return c.client.All(ctx)
}
