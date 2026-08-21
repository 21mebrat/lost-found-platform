package otp

import "errors"

var (
	ErrOTPSessionNotFound  = errors.New("otp session not found")
	ErrOTPExpired          = errors.New("otp has expired")
	ErrOTPAttemptsExceeded = errors.New("otp verification attempts exceeded")
	ErrInvalidOTP          = errors.New("invalid otp")
)
