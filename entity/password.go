package entity

import "fmt"

// Password is a password definition struct
type Password struct {
	UserID int
	Hash   string
	Salt   string
}

type PaswordRequest struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func (pr PaswordRequest) Validate() error {
	if len(pr.Old) < 7 {
		return fmt.Errorf("old password value must be longer than 7 characters, error: %w",
			ErrInvalidPassword)
	}

	if len(pr.New) < 7 {
		return fmt.Errorf("new password value must be longer than 7 characters, error: %w",
			ErrInvalidPassword)
	}

	if pr.Old == pr.New {
		return fmt.Errorf("new password must not be equal to old password error: %w",
			ErrInvalidPassword)
	}

	return nil
}
