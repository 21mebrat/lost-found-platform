package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrPhoneAlreadyUsed   = errors.New("phone already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserSuspended      = errors.New("user account suspended")
	ErrUserBlocked        = errors.New("user account blocked")
)
