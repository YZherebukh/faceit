package user

import (
	"bytes"
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
	updatePasswordURL = "http://localhost:8080/v1/users/%d/password"
)

var (
	testOldPassword = "old_password"
	testNewPassword = "new_password"
)

type testCaseUpdatePassword struct {
	url                string
	method             string
	input              entity.PaswordRequest
	expectedStatusCode int
}

func TestUpdatePassword(t *testing.T) {
	t.Run("positive_200", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusOK,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", 1)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)
		mockClientUpdatePassword.EXPECT().Update(ctx, testUserID, testNewPassword, testOldPassword).
			Return(nil)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_missing_userID", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_invalid_old_password", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: "short",
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_invalid_new_password", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: "short",
			},
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_equal_passwords", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testOldPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_client_validation_failed", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusBadRequest,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)
		mockClientUpdatePassword.EXPECT().Update(ctx, testUserID, testNewPassword, testOldPassword).
			Return(entity.ErrInvalidPassword)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_400_client_validation_failed", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusNotFound,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)
		mockClientUpdatePassword.EXPECT().Update(ctx, testUserID, testNewPassword, testOldPassword).
			Return(entity.ErrNotFound)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})

	t.Run("negative_500_client_error", func(t *testing.T) {
		tc := testCaseUpdatePassword{
			url:    fmt.Sprintf(updatePasswordURL, testUserID),
			method: http.MethodPut,
			input: entity.PaswordRequest{
				Old: testOldPassword,
				New: testNewPassword,
			},
			expectedStatusCode: http.StatusInternalServerError,
		}

		ctr := gomock.NewController(t)
		ctx := context.WithValue(context.Background(), "id", testUserID)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		logger := logger.New(mockLogger)

		mockClientUpdatePassword := mock_user.NewMockupdatePassword(ctr)
		mockClientUpdatePassword.EXPECT().Update(ctx, testUserID, testNewPassword, testOldPassword).
			Return(errTest)

		b, err := json.Marshal(tc.input)
		assert.Nil(t, err)

		req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(b)).WithContext(ctx)

		w := httptest.NewRecorder()

		newUpdatePassword(web.NewResponse(w, logger), mockClientUpdatePassword).Do(web.NewRequest(req))

		assert.Equal(t, tc.expectedStatusCode, w.Code)
	})
}
