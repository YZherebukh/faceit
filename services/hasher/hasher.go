package hasher

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	salt   string
	hashed []byte
}

func New() *Hasher {
	return &Hasher{
		salt: strconv.Itoa(bcrypt.MinCost),
	}
}

func (h *Hasher) HashAndSalt(password string) error {
	saltInt, err := strconv.Atoi(h.salt)
	if err != nil {
		return fmt.Errorf("unsupported salt format, %s, %w", h.salt, err)
	}

	h.hashed, err = bcrypt.GenerateFromPassword([]byte(password), saltInt)
	if err != nil {
		return fmt.Errorf("hashing password failed, error: %w", err)
	}

	return nil
}

func (h *Hasher) Salt() string {
	return h.salt
}

func (h *Hasher) Hashed() string {
	return string(h.hashed)
}
