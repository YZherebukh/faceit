package password

import (
	"context"
	"fmt"
	"testing"

	"github.com/faceit/test/entity"
	mock_password "github.com/faceit/test/services/password/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errTest = fmt.Errorf("errTest")

	testPasswordOne       = "qwerty"
	testPasswordTwo       = "zxcvbn"
	testPasswordHashedOne = "iweuriweurwueioru"
	testPasswordHashedTwo = "kenkekvclsmcmalsm"
	testSalt              = "test_salt"

	testUserID = 1
)

func TestUpdate(t *testing.T) {
	t.Run("positive_matching", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUpdate := mock_password.NewMockclient(ctr)
		mockUpdate.EXPECT().One(ctx, testUserID).
			Return(entity.Password{UserID: testUserID, Hash: testPasswordHashedOne}, nil)
		mockUpdate.EXPECT().Update(ctx, testUserID, testPasswordHashedTwo, testSalt).
			Return(nil)

		mockHasher := mock_password.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testPasswordOne, testPasswordHashedOne).Return(nil)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testPasswordTwo, testSalt).Return(testPasswordHashedTwo, nil)

		err := New(mockUpdate, mockHasher).Update(ctx, testUserID, testPasswordTwo, testPasswordOne)
		assert.Nil(t, err)
	})

	t.Run("positive_not_matching", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUpdate := mock_password.NewMockclient(ctr)
		mockUpdate.EXPECT().One(ctx, testUserID).
			Return(entity.Password{UserID: testUserID, Hash: testPasswordHashedOne}, nil)

		mockHasher := mock_password.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testPasswordOne, testPasswordHashedOne).Return(entity.ErrInvalidPassword)

		err := New(mockUpdate, mockHasher).Update(ctx, testUserID, testPasswordTwo, testPasswordOne)
		assert.ErrorIs(t, err, entity.ErrInvalidPassword)
	})

	t.Run("negative_failed_to_update", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUpdate := mock_password.NewMockclient(ctr)
		mockUpdate.EXPECT().One(ctx, testUserID).
			Return(entity.Password{UserID: testUserID, Hash: testPasswordHashedOne}, nil)
		mockUpdate.EXPECT().Update(ctx, testUserID, testPasswordHashedTwo, testSalt).
			Return(errTest)

		mockHasher := mock_password.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testPasswordOne, testPasswordHashedOne).Return(nil)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testPasswordTwo, testSalt).Return(testPasswordHashedTwo, nil)

		err := New(mockUpdate, mockHasher).Update(ctx, testUserID, testPasswordTwo, testPasswordOne)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("negative_failed_to_get_user", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUpdate := mock_password.NewMockclient(ctr)
		mockUpdate.EXPECT().One(ctx, testUserID).
			Return(entity.Password{UserID: testUserID, Hash: testPasswordHashedOne}, errTest)

		mockHasher := mock_password.NewMockhasher(ctr)

		err := New(mockUpdate, mockHasher).Update(ctx, testUserID, testPasswordTwo, testPasswordOne)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("negative_failed_to_hash", func(t *testing.T) {
		ctr := gomock.NewController(t)
		ctx := context.Background()

		mockUpdate := mock_password.NewMockclient(ctr)
		mockUpdate.EXPECT().One(ctx, testUserID).
			Return(entity.Password{UserID: testUserID, Hash: testPasswordHashedOne}, nil)

		mockHasher := mock_password.NewMockhasher(ctr)
		mockHasher.EXPECT().Compare(testPasswordOne, testPasswordHashedOne).Return(nil)
		mockHasher.EXPECT().Salt().Return(testSalt)
		mockHasher.EXPECT().Hash(testPasswordTwo, testSalt).Return("", errTest)

		err := New(mockUpdate, mockHasher).Update(ctx, testUserID, testPasswordTwo, testPasswordOne)
		assert.ErrorIs(t, err, errTest)
	})
}
