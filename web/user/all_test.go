package user

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
	"github.com/faceit/test/web"
	mock_user "github.com/faceit/test/web/user/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	allURL     = "http://localhost:8080/v1/users"
	testFilter = "filter_test"
	testTitle  = "title_test"
)

var (
	errTest = fmt.Errorf("errTest")

	testUser = entity.User{
		ID:        1,
		FirstName: "David",
		LastName:  "Bovie",
		NickName:  "Prince",
		Email:     "test@test.go",
		Password:  "qwerty",
		Salt:      "tre",
		Country:   "UK",
		CountryID: 1,
	}
)

type testCaseAll struct {
	url                string
	method             string
	filter             string
	title              string
	expectedResponse   []entity.UserResponse
	expectedStatusCode int
}

func (tc testCaseAll) checkresult(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, tc.expectedStatusCode, w.Code)

	if tc.expectedResponse != nil {
		var resp []entity.UserResponse

		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedResponse, resp)
	}
}

func TestAll(t *testing.T) {
	t.Run("positive_200_no_filters", func(t *testing.T) {
		tc := testCaseAll{
			url:                allURL,
			method:             http.MethodGet,
			expectedResponse:   []entity.UserResponse{testUser.ToResponse()},
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientAll := mock_user.NewMockall(ctr)
		mockClientAll.EXPECT().All(ctx, tc.title, tc.filter).Return([]entity.User{testUser}, nil)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)
		w := httptest.NewRecorder()

		newAll(web.NewResponse(w, logger), mockClientAll).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("positive_200_with_filters", func(t *testing.T) {
		tc := testCaseAll{
			url:                fmt.Sprintf("%s?title=%s&filter=%s", allURL, testTitle, testFilter),
			method:             http.MethodGet,
			filter:             testFilter,
			title:              testTitle,
			expectedResponse:   []entity.UserResponse{testUser.ToResponse()},
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientAll := mock_user.NewMockall(ctr)
		mockClientAll.EXPECT().All(ctx, tc.title, tc.filter).Return([]entity.User{testUser}, nil)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)
		w := httptest.NewRecorder()

		newAll(web.NewResponse(w, logger), mockClientAll).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("positive_204_no_users", func(t *testing.T) {
		tc := testCaseAll{
			url:                fmt.Sprintf("%s?title=%s&filter=%s", allURL, testTitle, testFilter),
			method:             http.MethodGet,
			filter:             testFilter,
			title:              testTitle,
			expectedStatusCode: http.StatusNoContent,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientAll := mock_user.NewMockall(ctr)
		mockClientAll.EXPECT().All(ctx, tc.title, tc.filter).Return(nil, entity.ErrNotFound)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)
		w := httptest.NewRecorder()

		newAll(web.NewResponse(w, logger), mockClientAll).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_500_client_error", func(t *testing.T) {
		tc := testCaseAll{
			url:                fmt.Sprintf("%s?title=%s&filter=%s", allURL, testTitle, testFilter),
			method:             http.MethodGet,
			filter:             testFilter,
			title:              testTitle,
			expectedStatusCode: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientAll := mock_user.NewMockall(ctr)
		mockClientAll.EXPECT().All(ctx, tc.title, tc.filter).Return(nil, errTest)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)
		w := httptest.NewRecorder()

		newAll(web.NewResponse(w, logger), mockClientAll).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})
}
