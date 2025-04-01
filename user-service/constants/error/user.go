package error

import "errors"

var (
	ErrUserNotFound         = errors.New("User is not found")
	ErrPasswordIncorrect    = errors.New("Password incorrect")
	ErrUsernameExist        = errors.New("Username already exist")
	ErrEmailExist           = errors.New("Email already exist")
	ErrPasswordDoesNotMatch = errors.New("Password does not match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordIncorrect,
	ErrUsernameExist,
	ErrPasswordDoesNotMatch,
}
