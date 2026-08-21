package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/auth"
	cryptoutil "github.com/21mebrat/lost-found-platform/internal/crypto"
	otpdomain "github.com/21mebrat/lost-found-platform/internal/domain/otp"
	userdomain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/21mebrat/lost-found-platform/internal/validator"
	"github.com/google/uuid"
)

const (
	otpSessionTTL    = 10 * time.Minute
	verifiedTokenTTL = 15 * time.Minute
)

var ErrorPhoneNotVerified = errors.New("phone number not verified or registration token expired")

func hashOTPCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

func getOTPSessionID(phone string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte("otp:"+phone))
}

func getTokenSessionID(token string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte("token:"+token))
}

// Step 1: Send OTP to Phone Number
func (s *Service) SendOTP(
	ctx context.Context,
	req SendOTPRequest,
) (*SendOTPResponse, error) {
	if strings.TrimSpace(req.Phone) == "" {
		return nil, ErrorPhoneRequired
	}

	phone, err := validator.ValidatePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	// Check if user is already registered in PostgreSQL
	existingUser, err := s.repo.GetByPhone(ctx, phone)
	if err != nil && !errors.Is(err, ErrorUserNotFound) {
		return nil, err
	}

	isRegistered := existingUser != nil

	otpCode, err := cryptoutil.GenerateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("generate otp code: %w", err)
	}

	purpose := otpdomain.PurposeRegistration
	if isRegistered {
		purpose = otpdomain.PurposeLogin
	}

	session := &otpdomain.OTPSession{
		ID:        getOTPSessionID(phone),
		Phone:     phone,
		CodeHash:  hashOTPCode(otpCode),
		Purpose:   purpose,
		Attempts:  0,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(otpSessionTTL),
	}

	if err := s.otpRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("save otp session: %w", err)
	}

	if err := s.smsService.SendOTP(ctx, phone, otpCode); err != nil {
		return nil, fmt.Errorf("send sms otp: %w", err)
	}

	msg := "Verification OTP sent successfully."
	if !isRegistered {
		msg = "Verification OTP sent. New registration required after verification."
	}

	return &SendOTPResponse{
		Message:          msg,
		Phone:            phone,
		IsRegistered:     isRegistered,
		ExpiresInSeconds: int(otpSessionTTL.Seconds()),
	}, nil
}

// Step 2: Verify OTP Code (Telegram Style)
func (s *Service) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	phone, err := validator.ValidatePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	otpCode := strings.TrimSpace(req.OTPCode)
	if otpCode == "" {
		return nil, ErrorInvalidOTP
	}

	sessionID := getOTPSessionID(phone)
	session, err := s.otpRepo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrorOTPNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.otpRepo.Delete(ctx, sessionID)
		return nil, ErrorOTPNotFound
	}

	if session.Attempts >= 3 {
		_ = s.otpRepo.Delete(ctx, sessionID)
		return nil, ErrorTooManyOTPAttempts
	}

	if session.CodeHash != hashOTPCode(otpCode) {
		attempts, err := s.otpRepo.IncrementAttempts(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if attempts >= 3 {
			_ = s.otpRepo.Delete(ctx, sessionID)
			return nil, ErrorTooManyOTPAttempts
		}
		return nil, ErrorInvalidOTP
	}

	// Delete OTP session on successful verification
	_ = s.otpRepo.Delete(ctx, sessionID)

	// Check if user is registered in PostgreSQL
	existingUser, err := s.repo.GetByPhone(ctx, phone)
	if err != nil && !errors.Is(err, ErrorUserNotFound) {
		return nil, err
	}

	isRegistered := existingUser != nil

	// If existing user -> Login path
	if isRegistered {
		return &VerifyOTPResponse{
			Message:      "Phone verified successfully. Welcome back!",
			Phone:        phone,
			IsRegistered: true,
			User:         ToUserResponse(existingUser),
		}, nil
	}

	// If new user -> Issue Registration Token for Step 3
	regToken, err := cryptoutil.GenerateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("generate registration token: %w", err)
	}

	tokenSessionID := getTokenSessionID(regToken)
	tokenSession := &otpdomain.OTPSession{
		ID:        tokenSessionID,
		Phone:     phone,
		Purpose:   otpdomain.PurposeRegistration,
		Attempts:  0,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(verifiedTokenTTL),
	}

	if err := s.otpRepo.Create(ctx, tokenSession); err != nil {
		return nil, fmt.Errorf("save verified token: %w", err)
	}

	return &VerifyOTPResponse{
		Message:           "Phone verified successfully. Please complete your registration profile.",
		Phone:             phone,
		IsRegistered:      false,
		RegistrationToken: regToken,
	}, nil
}

// Step 3: Complete Profile (Telegram Style - ONLY after phone is verified via registration_token)
func (s *Service) CompleteProfile(ctx context.Context, req CompleteProfileRequest) (*UserResponse, error) {
	regToken := strings.TrimSpace(req.RegistrationToken)
	if regToken == "" {
		return nil, ErrorPhoneNotVerified
	}

	// Verify registration token in Redis using the tokenSessionID
	tokenSessionID := getTokenSessionID(regToken)
	tokenSession, err := s.otpRepo.Get(ctx, tokenSessionID)
	if err != nil {
		return nil, ErrorPhoneNotVerified
	}
	if tokenSession == nil {
		return nil, ErrorPhoneNotVerified
	}

	if time.Now().After(tokenSession.ExpiresAt) {
		_ = s.otpRepo.Delete(ctx, tokenSessionID)
		return nil, ErrorPhoneNotVerified
	}

	reqPhone, err := validator.ValidatePhone(req.Phone)
	if err != nil || reqPhone != tokenSession.Phone {
		return nil, ErrorPhoneNotVerified
	}

	// Check if already registered in Postgres
	existingUser, err := s.repo.GetByPhone(ctx, tokenSession.Phone)
	if err != nil && !errors.Is(err, ErrorUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrorPhoneAlreadyExists
	}

	// Validate names
	firstName := strings.TrimSpace(req.FirstName)
	if firstName == "" || len(firstName) < 2 {
		return nil, ErrorFirstNameRequired
	}

	middleName := strings.TrimSpace(req.MiddleName)
	if middleName == "" || len(middleName) < 2 {
		return nil, ErrorMiddleNameRequired
	}

	lastName := strings.TrimSpace(req.LastName)
	if lastName == "" || len(lastName) < 2 {
		return nil, ErrorLastNameRequired
	}

	// Validate Email (optional)
	var emailPtr *string
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		cleanEmail := strings.ToLower(strings.TrimSpace(*req.Email))
		_, err := mail.ParseAddress(cleanEmail)
		if err != nil {
			return nil, ErrorInvalidEmail
		}

		existingByEmail, err := s.repo.GetByEmail(ctx, cleanEmail)
		if err != nil && !errors.Is(err, ErrorUserNotFound) {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, ErrorEmailAlreadyExists
		}
		emailPtr = &cleanEmail
	}

	// Validate Fayda (optional)
	var faydaPtr *string
	if req.Fayda != nil && strings.TrimSpace(*req.Fayda) != "" {
		cleanFayda := strings.TrimSpace(*req.Fayda)
		existingByFayda, err := s.repo.GetByFayda(ctx, cleanFayda)
		if err != nil && !errors.Is(err, ErrorUserNotFound) {
			return nil, err
		}
		if existingByFayda != nil {
			return nil, ErrorFaydaAlreadyExists
		}
		faydaPtr = &cleanFayda
	}

	// Language Preference
	lang := strings.ToLower(strings.TrimSpace(req.LanguageCode))
	if lang == "" {
		lang = "en"
	}

	// Validate password
	if strings.TrimSpace(req.Password) == "" {
		return nil, ErrorPasswordRequired
	}

	if !validator.ValidatePassword(req.Password) {
		return nil, ErrorInvalidPassword
	}

	hashedPassword, err := auth.PasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	// Insert into PostgreSQL with PhoneVerified = true
	newUser := &userdomain.User{
		FirstName:     firstName,
		MiddleName:    middleName,
		LastName:      lastName,
		Phone:         tokenSession.Phone,
		Email:         emailPtr,
		Fayda:         faydaPtr,
		LanguageCode:  lang,
		PasswordHash:  hashedPassword,
		PhoneVerified: true,
		Status:        userdomain.StatusActive,
	}

	createdUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user in db: %w", err)
	}

	// Consume and delete the registration token from Redis
	_ = s.otpRepo.Delete(ctx, tokenSessionID)

	return ToUserResponse(createdUser), nil
}
