package user

import (
	"bytes"
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
	createURL = "http://localhost:8080/v1/users"
)

var (
	testUserID    = 1
	testConsumers = []string{}
)

type testCaseCreate struct {
	url                string
	method             string
	consumers          []string
	input              entity.UserRequest
	expectedResponse   *createResponse
	expectedStatusCode int
}

type createResponse struct {
	ID int `json:"id"`
}

func (tc testCaseCreate) checkresult(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, tc.expectedStatusCode, w.Code)

	if tc.expectedResponse != nil {
		var resp createResponse

		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.Equal(t, *tc.expectedResponse, resp)
	}
}

func TestCreate(t *testing.T) {
	t.Run("positive_200", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwertyui",
				CountryID: 1,
			},
			expectedResponse:   &createResponse{ID: 1},
			consumers:          testConsumers,
			expectedStatusCode: http.StatusCreated,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger := logger.New(mockLogger)

		mockClientCreate := mock_user.NewMockcreate(ctr)
		mockClientCreate.EXPECT().Create(ctx, tc.input.ToUser()).Return(testUserID, nil)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		userNotify := tc.input
		userNotify.ID = testUserID
		mockNotifier.EXPECT().Add(entity.NotifierMessage{
			Message: entity.UserNotification{
				User:   userNotify.ToUser().ToResponse(),
				Action: actionCreate,
			},
			Consumers: tc.consumers,
		}).Times(1)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_400_missing_Firstname", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwertyui",
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

		mockClientCreate := mock_user.NewMockcreate(ctr)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_400_missing_LastName", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				FirstName: "David",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwertyui",
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

		mockClientCreate := mock_user.NewMockcreate(ctr)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_400_missing_NickName", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				Email:     "test@test.go",
				Password:  "qwertyui",
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

		mockClientCreate := mock_user.NewMockcreate(ctr)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_400_missing_email", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Password:  "qwertyui",
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

		mockClientCreate := mock_user.NewMockcreate(ctr)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})

	t.Run("negative_500_create_client_error", func(t *testing.T) {
		tc := testCaseCreate{
			url:    createURL,
			method: http.MethodPost,
			input: entity.UserRequest{
				FirstName: "David",
				LastName:  "Bovie",
				NickName:  "Prince",
				Email:     "test@test.go",
				Password:  "qwertyui",
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

		mockClientCreate := mock_user.NewMockcreate(ctr)
		mockClientCreate.EXPECT().Create(ctx, tc.input.ToUser()).Return(0, errTest)

		mockNotifier := mock_user.NewMocknotifier(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)
		w := httptest.NewRecorder()

		newCreate(web.NewResponse(w, logger), mockClientCreate, mockNotifier, testConsumers).Do(web.NewRequest(req))

		tc.checkresult(t, w)
	})
}
