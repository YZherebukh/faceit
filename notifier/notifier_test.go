package notifier

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/faceit/test/config"
	"github.com/faceit/test/entity"
	"github.com/faceit/test/logger"
	mock_logger "github.com/faceit/test/logger/mock"
	notifier_mock "github.com/faceit/test/notifier/mock"
	"github.com/golang/mock/gomock"
)

var (
	errTest = errors.New("error_test")

	testID      = 1
	testMessage = entity.NotifierMessage{Message: entity.UserResponse{
		ID: testID,
	},
		Consumers: testConsumers,
	}

	testMessageByte, _ = json.Marshal(testMessage)
	testConsumers      = []string{"test_consumer"}
	testConfig         = config.Notifier{
		Consumers:             config.Consumers{OnCreate: testConsumers},
		Timeout:               1,
		ClientMaxRetry:        3,
		ClientTimeoutIncrease: 1,
	}
)

func TestDo(t *testing.T) {
	t.Run("positive_one_try", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockNotifier := notifier_mock.NewMocknotifier(ctr)
		mockNotifier.EXPECT().Send(gomock.Any(), testConsumers, testMessageByte).Return(nil)

		mockLogger := mock_logger.NewMocklog(ctr)

		New(testConfig, mockNotifier, testConsumers, logger.New(mockLogger)).
			Do(ctx, testConsumers, testMessage)

	})

	t.Run("positive_2_tries", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockNotifier := notifier_mock.NewMocknotifier(ctr)
		mockNotifier.EXPECT().Send(gomock.Any(), testConsumers, testMessageByte).Return(errTest)
		mockNotifier.EXPECT().Send(gomock.Any(), testConsumers, testMessageByte).Return(nil)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).Times(1)

		New(testConfig, mockNotifier, testConsumers, logger.New(mockLogger)).
			Do(ctx, testConsumers, testMessage)
	})

	t.Run("positive_noConsumers", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockNotifier := notifier_mock.NewMocknotifier(ctr)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).Times(1)

		New(testConfig, mockNotifier, testConsumers, logger.New(mockLogger)).
			Do(ctx, []string{}, testMessage)

	})

	t.Run("negative", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockNotifier := notifier_mock.NewMocknotifier(ctr)
		mockNotifier.EXPECT().Send(gomock.Any(), testConsumers, testMessageByte).Return(errTest).Times(3)

		mockLogger := mock_logger.NewMocklog(ctr)
		mockLogger.EXPECT().Warningf(gomock.Any(), gomock.Any()).Times(3)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		New(testConfig, mockNotifier, testConsumers, logger.New(mockLogger)).
			Do(ctx, testConsumers, testMessage)
	})
}
