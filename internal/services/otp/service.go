package otp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/crypto"
	domain "github.com/21mebrat/lost-found-platform/internal/domain/otp"
	otprepo "github.com/21mebrat/lost-found-platform/internal/repository/otp"
	validator "github.com/21mebrat/lost-found-platform/internal/validator"
	"github.com/google/uuid"
)

const (
	otpLength = 6

	otpTTL = 5 * time.Minute

	maxAttempts = 5
)

type Service struct {
	repository otprepo.Repository
	hasher     Hasher
}

func NewService(
	repository otprepo.Repository,
	hasher Hasher,
) *Service {
	return &Service{
		repository: repository,
		hasher:     hasher,
	}
}

type CreateResult struct {
	Session *domain.OTPSession
	Code    string
}

func (s *Service) Create(
	ctx context.Context,
	phone string,
	purpose domain.Purpose,
) (*CreateResult, error) {

	phone = strings.TrimSpace(phone)

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	phone, err := validator.ValidatePhone(phone)

	if err != nil || phone == "" {
		return nil, errors.New("invalid phone number.")
	}

	if !validator.IsValidPurpose(purpose) {
		return nil, domain.ErrInvalidOTP
	}

	code, err := crypto.GenerateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	codeHash, err := s.hasher.Hash(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("hash otp: %w", err)
	}

	now := time.Now().UTC()

	session := &domain.OTPSession{
		ID:        uuid.New(),
		Phone:     phone,
		CodeHash:  codeHash,
		Purpose:   purpose,
		Attempts:  0,
		CreatedAt: now,
		ExpiresAt: now.Add(otpTTL),
	}

	if err := s.repository.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create otp session: %w", err)
	}

	return &CreateResult{
		Session: session,
		Code:    code,
	}, nil
}

func (s *Service) Verify(
	ctx context.Context,
	sessionID uuid.UUID,
	code string,
) error {

	code = strings.TrimSpace(code)

	if code == "" {
		return domain.ErrInvalidOTP
	}

	session, err := s.repository.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("get otp session: %w", err)
	}

	if session == nil {
		return domain.ErrOTPSessionNotFound
	}

	now := time.Now().UTC()

	if now.After(session.ExpiresAt) {
		_ = s.repository.Delete(ctx, session.ID)

		return domain.ErrOTPExpired
	}

	if session.Attempts >= maxAttempts {
		_ = s.repository.Delete(ctx, session.ID)

		return domain.ErrOTPAttemptsExceeded
	}

	valid, err := s.hasher.Compare(
		ctx,
		session.CodeHash,
		code,
	)
	if err != nil {
		return fmt.Errorf("compare otp: %w", err)
	}

	if !valid {
		attempts, err := s.repository.IncrementAttempts(
			ctx,
			session.ID,
		)
		if err != nil {
			return fmt.Errorf(
				"increment otp attempts: %w",
				err,
			)
		}

		if attempts >= maxAttempts {
			_ = s.repository.Delete(ctx, session.ID)

			return domain.ErrOTPAttemptsExceeded
		}

		return domain.ErrInvalidOTP
	}

	if err := s.repository.Delete(
		ctx,
		session.ID,
	); err != nil {
		return fmt.Errorf(
			"consume otp session: %w",
			err,
		)
	}

	return nil
}
