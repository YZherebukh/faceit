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
	"github.com/faceit/test/web/middleware"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
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
)

type testCaseAll struct {
	expectedBody   []entity.Country
	expectedStatus int
}

func TestAll(t *testing.T) {

	t.Run("positive_return_2_countries", func(t *testing.T) {
		tc := testCaseAll{
			expectedBody:   []entity.Country{testCountryUKR, testCountryUS},
			expectedStatus: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		logger := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)
		middlevare := middleware.New(logger)

		router := mux.NewRouter().StrictSlash(true)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockCountryClient.EXPECT().All(ctx).Return(tc.expectedBody, nil).Times(1)
		// mockLogger.EXPECT().Infof().Do()

		h := &Handler{router, logger, middlevare, countryService}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil).WithContext(ctx)

		w := httptest.NewRecorder()
		h.All(w, req)

		tc.checkAllResults(t, w)
	})

	t.Run("positive_no_countries", func(t *testing.T) {
		tc := testCaseAll{
			expectedStatus: http.StatusNoContent,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		logger := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)
		middlevare := middleware.New(logger)

		router := mux.NewRouter().StrictSlash(true)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockCountryClient.EXPECT().All(ctx).Return(nil, entity.ErrNotFound).Times(1)

		h := &Handler{router, logger, middlevare, countryService}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil).WithContext(ctx)

		w := httptest.NewRecorder()
		h.All(w, req)

		tc.checkAllResults(t, w)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		tc := testCaseAll{
			expectedStatus: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		logger := logger.New(mockLogger)
		mockCountryClient := mock_country.NewMockclient(ctr)
		countryService := country.New(mockCountryClient)
		middlevare := middleware.New(logger)

		router := mux.NewRouter().StrictSlash(true)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		mockCountryClient.EXPECT().All(ctx).Return([]entity.Country{}, errTest).Times(1)

		h := &Handler{router, logger, middlevare, countryService}

		req := httptest.NewRequest("GET", "http://example.com/foo", nil).WithContext(ctx)

		w := httptest.NewRecorder()
		h.All(w, req)

		tc.checkAllResults(t, w)
	})
}

func (tc testCaseAll) checkAllResults(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, tc.expectedStatus, w.Code)

	if tc.expectedBody != nil {
		var body []entity.Country

		err := json.NewDecoder(w.Body).Decode(&body)
		if err != nil {
			t.Errorf("failed to unmarshal body, error: %s", err)
		}

		assert.Equal(t, tc.expectedBody, body)
	}
}
