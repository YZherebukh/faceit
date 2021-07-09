package entity

import (
	"fmt"
	"regexp"
)

// User is a user definition struct
type User struct {
	ID        int
	FirstName string
	LastName  string
	NickName  string
	Email     string
	Password  string
	Salt      string
	Country   string
	CountryID int
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		NickName:  u.NickName,
		Email:     u.Email,
		Country:   u.Country,
	}
}

type UserResponse struct {
	ID        int    `jspn:"id"`
	FirstName string `jspn:"firstNmae"`
	LastName  string `jspn:"lastName"`
	NickName  string `jspn:"nickName"`
	Email     string `jspn:"email"`
	Country   string `jspn:"country"`
}

type UserNotification struct {
	User   UserResponse
	Action string `json:"action"`
}

type UserRequest struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstNmae"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CountryID int    `json:"country"`
}

func (ur UserRequest) ToUser() User {
	return User{
		ID:        ur.ID,
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		NickName:  ur.NickName,
		Password:  ur.Password,
		Email:     ur.Email,
		CountryID: ur.CountryID,
	}
}
func (u UserRequest) Validate() error {
	switch "" {
	case u.FirstName:
		return fmt.Errorf("%w, first name must not be empty", ErrValidationFailed)
	case u.LastName:
		return fmt.Errorf("%w, last name must not be empty", ErrValidationFailed)
	case u.NickName:
		return fmt.Errorf("%w, nick name must not be empty", ErrValidationFailed)
	case u.Password:
		return fmt.Errorf("%w, password must not be empty", ErrValidationFailed)
	case u.Email:
		return fmt.Errorf("%w, email name must not be empty", ErrValidationFailed)
	}

	pattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]" +
		"{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	regexp, err := regexp.Compile(pattern)
	if err != nil || !regexp.MatchString(u.Email) {
		return fmt.Errorf("invalid email format, %w", ErrValidationFailed)
	}

	return nil
}
