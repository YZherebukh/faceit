package user

import (
	"context"
	"errors"
	"testing"

	"github.com/faceit/test/entity"
	mock_user "github.com/faceit/test/services/user/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errTest = errors.New("error_test")

	testUserID               = 1
	testFirstName            = "David"
	testLastName             = "Bovie"
	testNickName             = "Prince"
	testEmail                = "imAwesome@everything.go"
	testPassword             = "qwerty"
	testSalt                 = "test_salt"
	testCountryID            = 100
	testCountryName          = "UK"
	testFilterParamCountry   = "country"
	testFilterParamFirstName = "firstName"
	testFilterParamLastName  = "lastName"
	testFilterParamNickName  = "nickName"
	testFilterParamEmail     = "email"

	testFilterFirstName = "first_name"
	testFilterLastName  = "last_name"
	testFilterNickName  = "nick_name"
	testFilterEmail     = "email"

	testPasswordHased = "efbwebvjbjwqencj"
	testUser          = entity.User{
		FirstName: testFirstName,
		LastName:  testLastName,
		NickName:  testNickName,
		Email:     testEmail,
		Password:  testPassword,
		CountryID: testCountryID,
	}

	testUserupdate = entity.User{
		ID:        testUserID,
		FirstName: testFirstName,
		LastName:  testLastName,
		NickName:  "Freddy",
		Email:     testEmail,
		Password:  testPassword,
		CountryID: testCountryID,
	}

	testUserHashedPassword = entity.User{
		FirstName: testFirstName,
		LastName:  testLastName,
		NickName:  testNickName,
		Email:     testEmail,
		Password:  testPasswordHased,
		Salt:      testSalt,
		CountryID: testCountryID,
	}

	testPasswordStr = entity.Password{
		Hash: testPasswordHased,
	}

	testUsers = []entity.User{testUserHashedPassword}
)

func TestCreate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Create(ctx, testUserHashedPassword).Return(testUserID, nil)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testUser.Password, testSalt).Return(testPasswordHased, nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		id, err := New(mockUserClient, mockHasher, mockPassword).Create(ctx, testUser)
		assert.Nil(t, err)
		assert.Equal(t, testUserID, id)
	})

	t.Run("negative_failed_to_salt", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testUser.Password, testSalt).Return("", errTest)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		id, err := New(mockUserClient, mockHasher, mockPassword).Create(ctx, testUser)
		assert.ErrorIs(t, err, errTest)
		assert.Equal(t, 0, id)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Create(ctx, testUserHashedPassword).Return(0, errTest)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testUser.Password, testSalt).Return(testPasswordHased, nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		id, err := New(mockUserClient, mockHasher, mockPassword).Create(ctx, testUser)
		assert.ErrorIs(t, err, errTest)
		assert.Equal(t, 0, id)
	})
}

func TestAll(t *testing.T) {
	t.Run("positive_all_by_country", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().AllByCountry(ctx, testCountryName).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, testFilterParamCountry, testCountryName)
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("positive_all_by_first_name", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().AllWithFilter(ctx, testFilterFirstName, testFirstName).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, testFilterParamFirstName, testFirstName)
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("positive_all_by_last_name", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().AllWithFilter(ctx, testFilterLastName, testLastName).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, testFilterParamLastName, testLastName)
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("positive_all_by_nick_name", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().AllWithFilter(ctx, testFilterNickName, testNickName).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, testFilterParamNickName, testNickName)
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("positive_all_by_email", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().AllWithFilter(ctx, testFilterEmail, testEmail).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, testFilterParamEmail, testEmail)
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("positive_all_no_filters", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().All(ctx).Return(testUsers, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, "", "")
		assert.Nil(t, err)
		assert.Equal(t, testUsers, countries)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().All(ctx).Return(nil, errTest)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		countries, err := New(mockUserClient, mockHasher, mockPassword).All(ctx, "", "")
		assert.Nil(t, countries)
		assert.ErrorIs(t, err, errTest)
	})
}

func TestOne(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().One(ctx, testUserID).Return(testUser, nil)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		user, err := New(mockUserClient, mockHasher, mockPassword).One(ctx, testUserID)
		assert.Nil(t, err)
		assert.Equal(t, testUser, user)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().One(ctx, testUserID).Return(entity.User{}, errTest)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)

		user, err := New(mockUserClient, mockHasher, mockPassword).One(ctx, testUserID)
		assert.Equal(t, user, entity.User{})
		assert.ErrorIs(t, err, errTest)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Update(ctx, testUserupdate).Return(nil)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Update(ctx, testUserupdate)
		assert.Nil(t, err)
	})

	t.Run("negative_validation_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(entity.Password{}, errTest)

		err := New(mockUserClient, mockHasher, mockPassword).Update(ctx, testUserupdate)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("negative_invalid_password", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(entity.ErrInvalidPassword)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Update(ctx, testUserupdate)
		assert.ErrorIs(t, err, entity.ErrInvalidPassword)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Update(ctx, testUserupdate).Return(errTest)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Update(ctx, testUserupdate)
		assert.ErrorIs(t, err, errTest)
	})
}

func TestDelete(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Delete(ctx, testUserupdate.ID).Return(nil)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Delete(ctx, testUserupdate)
		assert.Nil(t, err)
	})

	t.Run("negative_validation_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(entity.Password{}, errTest)

		err := New(mockUserClient, mockHasher, mockPassword).Delete(ctx, testUserupdate)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("negative_invalid_password", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(entity.ErrInvalidPassword)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Delete(ctx, testUserupdate)
		assert.ErrorIs(t, err, entity.ErrInvalidPassword)
	})

	t.Run("negative_client_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)
		mockUserClient.EXPECT().Delete(ctx, testUserupdate.ID).Return(errTest)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).Delete(ctx, testUserupdate)
		assert.ErrorIs(t, err, errTest)
	})
}

func TestCanUpdate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(nil)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).canUpdate(ctx, testUserupdate)
		assert.Nil(t, err)
	})

	t.Run("negative_validation_error", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(entity.Password{}, errTest)

		err := New(mockUserClient, mockHasher, mockPassword).canUpdate(ctx, testUserupdate)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("negative_invalid_password", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUserClient := mock_user.NewMockclient(ctr)

		mockHasher := mock_user.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testUserupdate.Password, testUserHashedPassword.Password).Return(entity.ErrInvalidPassword)

		mockPassword := mock_user.NewMockpasswordClient(ctr)
		mockPassword.EXPECT().One(ctx, testUserID).Return(testPasswordStr, nil)

		err := New(mockUserClient, mockHasher, mockPassword).canUpdate(ctx, testUserupdate)
		assert.ErrorIs(t, err, entity.ErrInvalidPassword)
	})
}
