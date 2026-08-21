package user

import "errors"

var (
	ErrorFirstNameRequired  = errors.New("first name is required")
	ErrorMiddleNameRequired = errors.New("middle name is required")
	ErrorLastNameRequired   = errors.New("last name is required")
	ErrorEmailRequired      = errors.New("email is required")
	ErrorInvalidEmail       = errors.New("invalid email address format")
	ErrorPhoneRequired      = errors.New("phone number is required")
	ErrorPasswordRequired   = errors.New("password is required")
	ErrorInvalidPassword    = errors.New("password must be at least 8 characters and contain uppercase, lowercase, digit, and special character")
	ErrorPhoneAlreadyExists = errors.New("phone number already registered")
	ErrorEmailAlreadyExists = errors.New("email address already registered")
	ErrorFaydaAlreadyExists = errors.New("fayda ID already registered")
	ErrorUserNotFound       = errors.New("user not found")
	ErrorInvalidID          = errors.New("invalid user ID format")
	ErrorOTPNotFound        = errors.New("otp expired or registration session not found")
	ErrorInvalidOTP         = errors.New("invalid otp verification code")
	ErrorTooManyOTPAttempts = errors.New("too many failed otp attempts, please request a new code")
	ErrorInvalidIput        = "invalid request body"
)
