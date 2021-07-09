package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/faceit/test/entity"
)

// countries table parameters and query
const (
	countryTable  = `countries`
	countryParams = `country_id, iso2, country_name`

	selectOneCountryQuery   = `SELECT ` + countryParams + ` FROM ` + countryTable + ` WHERE country_id = $1;`
	selectAllCountriesQuery = `SELECT ` + countryParams + ` FROM ` + countryTable + `;`
)

// Country is a country store implementation
type Country struct {
	*sql.DB
}

// NewCountry creates a Country instance
func NewCountry(db *sql.DB) *Country {
	return &Country{
		db,
	}
}

// All returns list with all countries
func (c *Country) All(ctx context.Context) ([]entity.Country, error) {
	rows, err := c.QueryContext(ctx, selectAllCountriesQuery)
	if err != nil {
		return nil, fmt.Errorf("query failed, %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	countries := []entity.Country{}

	for rows.Next() {
		country := entity.Country{}

		err = rows.Scan(&country.ID, &country.ISO2, &country.Name)
		if err != nil {
			return nil, fmt.Errorf("query failed, %w", err)
		}

		countries = append(countries, country)
	}

	return countries, nil
}

// One returns one country by it's id
func (c *Country) One(ctx context.Context, id int) (entity.Country, error) {
	country := entity.Country{}

	err := c.QueryRowContext(ctx, selectOneCountryQuery, id).Scan(
		&country.ID,
		&country.Name,
		&country.ISO2)
	if err != nil {
		err = fmt.Errorf("query failed, %w", err)
	}

	return country, err
}
