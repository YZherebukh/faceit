package country

import (
	"context"
	"fmt"
	"testing"

	"github.com/faceit/test/entity"
	mock_country "github.com/faceit/test/services/country/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errTest        = fmt.Errorf("errTest")
	testCountryUKR = entity.Country{
		ID:   1,
		Name: "Ukraine",
		ISO2: "UKR",
	}

	testCountryUS = entity.Country{
		ID:   2,
		Name: "USA",
		ISO2: "US",
	}

	testCountries = []entity.Country{testCountryUKR, testCountryUS}
)

func TestAll(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockAll := mock_country.NewMockclient(ctr)
		mockAll.EXPECT().All(ctx).Return(testCountries, nil)

		countries, err := New(mockAll).All(ctx)
		assert.Nil(t, err)
		assert.Equal(t, testCountries, countries)
	})

	t.Run("positive_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockAll := mock_country.NewMockclient(ctr)
		mockAll.EXPECT().All(ctx).Return(nil, errTest)

		countries, err := New(mockAll).All(ctx)
		assert.Nil(t, countries)
		assert.Equal(t, errTest, err)
	})
}
