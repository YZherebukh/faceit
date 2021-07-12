package user

import (
	"context"
	"encoding/json"
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
	oneURL = "http://localhost:8080/v1/users"
)

var (
	testUserResponse = testUser.ToResponse()
)

type testCaseOne struct {
	url                string
	method             string
	expectedResponse   *entity.UserResponse
	expectedStatusCode int
}

func (tc testCaseOne) checkresult(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, tc.expectedStatusCode, w.Code)

	if tc.expectedResponse != nil {
		var resp entity.UserResponse

		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.Equal(t, *tc.expectedResponse, resp)
	}
}

func TestOne(t *testing.T) {
	t.Run("positive_200", func(t *testing.T) {
		tc := testCaseOne{
			url:                oneURL,
			method:             http.MethodGet,
			expectedResponse:   &testUserResponse,
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientOne := mock_user.NewMockone(ctr)
		mockClientOne.EXPECT().One(ctx, testUserID).Return(testUser, nil)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()

		newOne(web.NewResponse(w, logger), mockClientOne).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_400_missing_userId", func(t *testing.T) {
		tc := testCaseOne{
			url:                oneURL,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientOne := mock_user.NewMockone(ctr)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()

		newOne(web.NewResponse(w, logger), mockClientOne).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_404_no_content", func(t *testing.T) {
		tc := testCaseOne{
			url:                oneURL,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusNotFound,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientOne := mock_user.NewMockone(ctr)
		mockClientOne.EXPECT().One(ctx, testUserID).Return(entity.User{}, entity.ErrNotFound)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()

		newOne(web.NewResponse(w, logger), mockClientOne).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_500_client_one_error", func(t *testing.T) {
		tc := testCaseOne{
			url:                oneURL,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientOne := mock_user.NewMockone(ctr)
		mockClientOne.EXPECT().One(ctx, testUserID).Return(entity.User{}, errTest)

		req := httptest.NewRequest(tc.method, tc.url, nil).WithContext(ctx)

		w := httptest.NewRecorder()

		newOne(web.NewResponse(w, logger), mockClientOne).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})
}
