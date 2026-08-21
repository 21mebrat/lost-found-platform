package otp

import (
	"time"

	"github.com/google/uuid"
)

type Purpose string

const (
	PurposeLogin         Purpose = "login"
	PurposeRegistration  Purpose = "registration"
	PurposePasswordReset Purpose = "password_reset"
	PurposePhoneChange   Purpose = "phone_change"
)

type OTPSession struct {
	ID        uuid.UUID
	Phone     string
	CodeHash  string
	Purpose   Purpose
	Attempts  int
	CreatedAt time.Time
	ExpiresAt time.Time
}
