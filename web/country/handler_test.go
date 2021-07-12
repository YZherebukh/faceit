package country

import (
	"testing"

	"github.com/faceit/test/logger"
	mock_logger "github.com/faceit/test/logger/mock"
	"github.com/faceit/test/services/country"
	mock_country "github.com/faceit/test/services/country/mock"
	"github.com/faceit/test/web/middleware"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestHandler(t *testing.T) {
	t.Run("positive_return_2_countries", func(t *testing.T) {
		ctr := gomock.NewController(t)

		mockLogger := mock_logger.NewMocklog(ctr)
		log := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)

		NewHandler(mux.NewRouter().StrictSlash(true), log, middleware.New(log), countryService)
	})
}
