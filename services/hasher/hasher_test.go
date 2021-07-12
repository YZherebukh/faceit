package hasher

import (
	"testing"

	"github.com/faceit/test/entity"
	"github.com/stretchr/testify/assert"
)

var (
	testPasswordOne = "qwerty"
	testPasswordTwo = "zxcvbn"
)

func TestCompare(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		hash := New()

		pass, err := hash.Hash(testPasswordOne, hash.Salt())
		assert.Nil(t, err)

		err = hash.Compare(testPasswordOne, pass)
		assert.Nil(t, err)
	})

	t.Run("not_matching", func(t *testing.T) {
		hash := New()

		pass, err := hash.Hash(testPasswordOne, hash.Salt())
		assert.Nil(t, err)

		err = hash.Compare(testPasswordTwo, pass)
		assert.ErrorIs(t, err, entity.ErrInvalidPassword)
	})
}
