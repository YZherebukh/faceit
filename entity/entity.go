package entity

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrValidationFailed = errors.New("validation failed")
	ErrInvalidFilter    = errors.New("invalid filter")
	ErrInvalidPassword  = errors.New("invalid filter")
	ErrUserExist        = errors.New("user exist")
	ErrUserDoesNotExist = errors.New("user does not exist")
	ErrUserIDIsMissing  = errors.New("user id is missing")
)
