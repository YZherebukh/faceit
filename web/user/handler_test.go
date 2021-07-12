package user

import (
	"testing"

	"github.com/faceit/test/config"
	"github.com/faceit/test/logger"
	mock_logger "github.com/faceit/test/logger/mock"
	"github.com/faceit/test/queue"
	queue_mock "github.com/faceit/test/queue/mock"
	"github.com/faceit/test/services/country"
	mock_country "github.com/faceit/test/services/country/mock"
	"github.com/faceit/test/services/hasher"
	"github.com/faceit/test/services/password"
	mock_password "github.com/faceit/test/services/password/mock"
	"github.com/faceit/test/services/user"
	mock_user "github.com/faceit/test/services/user/mock"
	"github.com/faceit/test/web/middleware"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestNewhandler(t *testing.T) {
	ctr := gomock.NewController(t)

	mockLogger := mock_logger.NewMocklog(ctr)

	logger := logger.New(mockLogger)

	mockCountryClient := mock_country.NewMockclient(ctr)
	mockCountry := country.New(mockCountryClient)

	mockPasswordClient := mock_password.NewMockclient(ctr)
	mockPasswordHasher := mock_password.NewMockhasher(ctr)
	mockPassword := password.New(mockPasswordClient, mockPasswordHasher)

	mockUserClient := mock_user.NewMockclient(ctr)
	mockUserHasher := mock_user.NewMockhasher(ctr)
	mockUser := user.New(mockUserClient, mockUserHasher, mockPasswordClient)

	hasher := hasher.New()
	mockNotifier := queue_mock.NewMocknotifier(ctr)
	queue := queue.New(config.Queue{}, mockNotifier)

	NewHandler(
		mux.NewRouter().StrictSlash(true),
		logger,
		middleware.New(logger),
		mockUser,
		mockCountry,
		mockPassword,
		hasher,
		*queue,
	)
}
