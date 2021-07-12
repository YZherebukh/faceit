package hasher

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/faceit/test/entity"
	"golang.org/x/crypto/bcrypt"
)

// Hasher is a hasher struct
type Hasher struct {
}

// New create new Hasher instance
func New() *Hasher {
	return &Hasher{}
}

// Hash salts passwrd with salt and hashes it
func (h *Hasher) Hash(password, salt string) (string, error) {
	saltInt, err := strconv.Atoi(salt)
	if err != nil {
		return "", fmt.Errorf("unsupported salt format, %s, %w", salt, err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), saltInt)
	if err != nil {
		return "", fmt.Errorf("hashing password failed, error: %w", err)
	}

	return string(hashed), nil
}

// Salt returnes hasher salt
func (h *Hasher) Salt() string {
	return strconv.Itoa(bcrypt.MinCost)
}

// Compare compares hashed and unhashed password
func (h *Hasher) Compare(password, hashed string) error {
	byteHash := []byte(hashed)

	err := bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return entity.ErrInvalidPassword
	}

	if err != nil {
		return fmt.Errorf("filed to compare passwords, error: %w", err)
	}

	return nil
}
