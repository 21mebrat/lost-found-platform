package user

import "errors"

var (
	ErrorFirstNameRequired = errors.New("First name is required.")
	ErrorLastNameRequired  = errors.New("Last name is required.")
	ErrorEmailRequired     = errors.New("email required.")
	ErrorInvalidEmail      = errors.New("Invalid email address")
	ErrorPhoneRequired     = errors.New("Phone number is required.")
	ErrorPasswordRequired  = errors.New("Password is required.")
	ErrorInvalidPassword   = errors.New(
		"password must be at least 8 characters and contain uppercase, lowercase, digit, and special character",
	)
	ErrorPhoneAlreadyExists = errors.New("The Phone Already exists.")
	ErrorEmailAlreadyExists = errors.New("The Email already exists.")
	ErrorUserNotFound       = errors.New("user is not found.")
	ErrorInvalidIput        = "invalid request body"
)
