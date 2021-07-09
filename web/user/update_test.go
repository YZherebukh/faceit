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
	updateURL = "http://localhost:8080/v1/users"
)

type testCaseUpdate struct {
	url                string
	urlParams          map[string][]string
	method             string
	consumers          []string
	input              entity.UserRequest
	expectedStatusCode int
}

func TestUpdate(t *testing.T) {
	t.Run("positive_200", func(t *testing.T) {
		tc := testCaseUpdate{
			url:       fmt.Sprintf("%s/%d", updateURL, testUserID),
			urlParams: map[string][]string{"id": {strconv.Itoa(testUserID)}},

			method: http.MethodPut,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwerty",
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

		mockClientUpdate := mock_user.NewMockupdate(ctr)
		userUpdate := tc.input.ToUser()
		userUpdate.ID = testUserID
		mockClientUpdate.EXPECT().Update(ctx, userUpdate).Return(nil)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		userNotify := tc.input
		userNotify.ID = testUserID
		mockNotifier.EXPECT().Add(entity.NotifierMessage{
			Message: entity.UserNotification{
				User:   userNotify.ToUser().ToResponse(),
				Action: actionUpdate,
			},
			Consumers: tc.consumers,
		}).Times(1)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams

		w := httptest.NewRecorder()

		newUpdate(web.NewResponse(w, logger), mockClientUpdate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_user_not_found", func(t *testing.T) {
		tc := testCaseUpdate{
			url:       fmt.Sprintf("%s/%d", updateURL, testUserID),
			urlParams: map[string][]string{"id": {strconv.Itoa(testUserID)}},

			method: http.MethodPut,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwerty",
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

		mockClientUpdate := mock_user.NewMockupdate(ctr)
		userUpdate := tc.input.ToUser()
		userUpdate.ID = testUserID
		mockClientUpdate.EXPECT().Update(ctx, userUpdate).Return(entity.ErrNotFound)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		req.PostForm = tc.urlParams

		w := httptest.NewRecorder()

		newUpdate(web.NewResponse(w, logger), mockClientUpdate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_missing_userID", func(t *testing.T) {
		tc := testCaseUpdate{
			url: fmt.Sprintf("%s/%d", updateURL, testUserID),

			method: http.MethodPut,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwerty",
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

		mockClientUpdate := mock_user.NewMockupdate(ctr)
		userUpdate := tc.input.ToUser()
		userUpdate.ID = testUserID

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdate(web.NewResponse(w, logger), mockClientUpdate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})
}
