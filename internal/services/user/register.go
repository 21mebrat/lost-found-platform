package user

import (
	"context"
	"net/mail"
	"regexp"
	"strings"

	"github.com/21mebrat/lost-found-platform/internal/auth"
	"github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/21mebrat/lost-found-platform/internal/validator"
)

var passwordRegex = regexp.MustCompile(
	`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]).{8,}$`,
)

func (s *Service) Register(ctx context.Context, input RegisterRequest) (*UserResponse, error) {
	if strings.TrimSpace(input.FirstName) == "" || len(strings.TrimSpace(input.FirstName)) < 2 {
		return nil, ErrorFirstNameRequired
	}

	if strings.TrimSpace(input.LastName) == "" || len(strings.TrimSpace(input.LastName)) < 2 {
		return nil, ErrorLastNameRequired
	}

	if strings.TrimSpace(input.Email) == "" {
		return nil, ErrorEmailRequired
	}

	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return nil, ErrorInvalidEmail
	}

	if strings.TrimSpace(input.Phone) == "" {
		return nil, ErrorPhoneRequired
	}

	phone, err := validator.ValidatePhone(input.Phone)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Password) == "" {
		return nil, ErrorPasswordRequired
	}

	if !passwordRegex.MatchString(input.Password) {
		return nil, ErrorInvalidPassword
	}

	hashedPassword, err := auth.PasswordHash(input.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		FirstName:    strings.TrimSpace(input.FirstName),
		LastName:     strings.TrimSpace(input.LastName),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		Phone:        phone,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:              createdUser.ID.String(),
		FirstName:       createdUser.FirstName,
		LastName:        createdUser.LastName,
		Email:           createdUser.Email,
		Phone:           createdUser.Phone,
		EmailVerified:   createdUser.EmailVerified,
		PhoneVerified:   createdUser.PhoneVerified,
		ProfileImageURL: createdUser.ProfileImageURL,
		Role:            string(createdUser.Role),
		Status:          string(createdUser.Status),
	}, nil
}
