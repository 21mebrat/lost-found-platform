package user

import (
	otprepo "github.com/21mebrat/lost-found-platform/internal/repository/otp"
	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
	smsservice "github.com/21mebrat/lost-found-platform/internal/services/sms"
)

// Service provides application business logic for User entity operations and OTP phone verification.
type Service struct {
	repo       userrepo.Repository
	otpRepo    otprepo.Repository
	smsService smsservice.SMSService
}

// Compile-time interface assertion
var _ UserService = (*Service)(nil)

// NewService constructs a new User service instance with required dependencies.
func NewService(
	repo userrepo.Repository,
	otpRepo otprepo.Repository,
	smsService smsservice.SMSService,
) *Service {
	return &Service{
		repo:       repo,
		otpRepo:    otpRepo,
		smsService: smsService,
	}
}
