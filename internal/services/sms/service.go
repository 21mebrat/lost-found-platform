package sms

import (
	"context"
	"fmt"

	"github.com/21mebrat/lost-found-platform/internal/notification"
)

type service struct {
	notificationService *notification.Service
}

func NewService(ns *notification.Service) SMSService {
	return &service{
		notificationService: ns,
	}
}

func (s *service) SendOTP(ctx context.Context, phone, code string) error {
	req := notification.Request{
		Channel:   notification.ChannelSMS,
		Recipient: phone,
		Message:   fmt.Sprintf("Your OTP verification code is: %s", code),
	}
	return s.notificationService.Send(ctx, req)
}
