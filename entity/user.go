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

// ToResponse is transforming User struct to UserResponse struct
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

// UserResponse is a user response struct
type UserResponse struct {
	ID        int    `jspn:"id"`
	FirstName string `jspn:"firstNmae"`
	LastName  string `jspn:"lastName"`
	NickName  string `jspn:"nickName"`
	Email     string `jspn:"email"`
	Country   string `jspn:"country"`
}

// UserNotification is a user notification struct
type UserNotification struct {
	User   UserResponse
	Action string `json:"action"`
}

// UserRequest is a user request struct
type UserRequest struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstNmae"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CountryID int    `json:"country"`
}

// ToUser transformes UserRequest struct to User struct
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

// Validate validates user password
func (u UserRequest) Validate() error {
	switch "" {
	case u.FirstName:
		return fmt.Errorf("%w, first name must not be empty", ErrValidationFailed)
	case u.LastName:
		return fmt.Errorf("%w, last name must not be empty", ErrValidationFailed)
	case u.NickName:
		return fmt.Errorf("%w, nick name must not be empty", ErrValidationFailed)
	case u.Email:
		return fmt.Errorf("%w, email name must not be empty", ErrValidationFailed)
	}

	if len(u.Password) < 7 {
		return fmt.Errorf("%w, password must be at least 7 charecters long", ErrValidationFailed)
	}

	pattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]" +
		"{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	regexp, err := regexp.Compile(pattern)
	if err != nil || !regexp.MatchString(u.Email) {
		return fmt.Errorf("invalid email format, %w", ErrValidationFailed)
	}

	return nil
}
