package sms

import "context"

type SMSService interface {
	SendOTP(ctx context.Context, phone, code string) error
}
