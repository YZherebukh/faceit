package country

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faceit/test/entity"
	"github.com/faceit/test/logger"
	mock_logger "github.com/faceit/test/logger/mock"
	"github.com/faceit/test/services/country"
	mock_country "github.com/faceit/test/services/country/mock"
	"github.com/faceit/test/web"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

const (
	allURL     = "http://localhost:8080/v1/countries"
	testFilter = "filter_test"
	testTitle  = "title_test"
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
)

type testCaseAll struct {
	url                string
	method             string
	expectedResponse   []entity.Country
	expectedStatusCode int
}

func (tc testCaseAll) checkAllResults(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, tc.expectedStatusCode, w.Code)

	if tc.expectedResponse != nil {
		var body []entity.Country

		err := json.NewDecoder(w.Body).Decode(&body)
		if err != nil {
			t.Errorf("failed to unmarshal body, error: %s", err)
		}

		assert.Equal(t, tc.expectedResponse, body)
	}
}

func TestAll(t *testing.T) {
	t.Run("positive_return_2_countries", func(t *testing.T) {
		tc := testCaseAll{
			url:                allURL,
			method:             http.MethodGet,
			expectedResponse:   []entity.Country{testCountryUKR, testCountryUS},
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		log := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)

		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockCountryClient.EXPECT().All(ctx).Return(tc.expectedResponse, nil).Times(1)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()
		newAll(web.NewResponse(w, log), countryService).Do(web.NewRequest(req))

		tc.checkAllResults(t, w)
	})

	t.Run("positive_no_countries", func(t *testing.T) {
		tc := testCaseAll{
			url:                allURL,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusNoContent,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		log := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)

		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockCountryClient.EXPECT().All(ctx).Return(nil, entity.ErrNotFound).Times(1)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()
		newAll(web.NewResponse(w, log), countryService).Do(web.NewRequest(req))

		tc.checkAllResults(t, w)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		tc := testCaseAll{
			url:                allURL,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		log := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)

		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		mockCountryClient.EXPECT().All(ctx).Return([]entity.Country{}, errTest).Times(1)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()
		newAll(web.NewResponse(w, log), countryService).Do(web.NewRequest(req))

		tc.checkAllResults(t, w)
	})
}
