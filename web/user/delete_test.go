package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	deleteURL = "http://localhost:8080/v1/users"
)

var (
	testPassword = "qwerty"
)

type testCaseDelete struct {
	url                string
	urlParams          map[string][]string
	method             string
	consumers          []string
	input              entity.UserRequest
	expectedStatusCode int
}

func TestDelete(t *testing.T) {
	t.Run("positive_200", func(t *testing.T) {
		tc := testCaseDelete{
			url:       fmt.Sprintf("%s/%d", deleteURL, testUserID),
			urlParams: map[string][]string{"id": {strconv.Itoa(testUserID)}},
			method:    http.MethodDelete,
			input: entity.UserRequest{
				Password:  testPassword,
				CountryID: 1,
			},
			consumers:          testConsumers,
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientDelete := mock_user.NewMockdelete(ctr)

		deleteUser := tc.input
		deleteUser.ID = testUserID

		mockClientDelete.EXPECT().Delete(ctx, deleteUser.ToUser()).Return(nil)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		userNotify := tc.input
		userNotify.ID = testUserID
		mockNotifier.EXPECT().Add(entity.NotifierMessage{
			Message: entity.UserNotification{
				User:   deleteUser.ToUser().ToResponse(),
				Action: actionDelete,
			},
			Consumers: tc.consumers,
		}).Times(1)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams
		w := httptest.NewRecorder()

		newDelete(web.NewResponse(w, logger), mockClientDelete, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("positive_404_user_not_found", func(t *testing.T) {
		tc := testCaseDelete{
			url:       fmt.Sprintf("%s/%d", deleteURL, testUserID),
			urlParams: map[string][]string{"id": {strconv.Itoa(testUserID)}},
			method:    http.MethodDelete,
			input: entity.UserRequest{
				Password:  testPassword,
				CountryID: 1,
			},
			consumers:          testConsumers,
			expectedStatusCode: http.StatusNotFound,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientDelete := mock_user.NewMockdelete(ctr)

		deleteUser := tc.input
		deleteUser.ID = testUserID

		mockClientDelete.EXPECT().Delete(ctx, deleteUser.ToUser()).Return(entity.ErrNotFound)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams
		w := httptest.NewRecorder()

		newDelete(web.NewResponse(w, logger), mockClientDelete, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("positive_400_missing_userID", func(t *testing.T) {
		tc := testCaseDelete{
			url:    fmt.Sprintf("%s/%d", deleteURL, testUserID),
			method: http.MethodDelete,
			input: entity.UserRequest{
				Password:  testPassword,
				CountryID: 1,
			},
			consumers:          testConsumers,
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientDelete := mock_user.NewMockdelete(ctr)

		deleteUser := tc.input
		deleteUser.ID = testUserID

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams
		w := httptest.NewRecorder()

		newDelete(web.NewResponse(w, logger), mockClientDelete, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_500", func(t *testing.T) {
		tc := testCaseDelete{
			url:       fmt.Sprintf("%s/%d", deleteURL, testUserID),
			urlParams: map[string][]string{"id": {strconv.Itoa(testUserID)}},
			method:    http.MethodDelete,
			input: entity.UserRequest{
				Password:  testPassword,
				CountryID: 1,
			},
			consumers:          testConsumers,
			expectedStatusCode: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientDelete := mock_user.NewMockdelete(ctr)

		deleteUser := tc.input
		deleteUser.ID = testUserID

		mockClientDelete.EXPECT().Delete(ctx, deleteUser.ToUser()).Return(errTest)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams
		w := httptest.NewRecorder()

		newDelete(web.NewResponse(w, logger), mockClientDelete, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})
}
