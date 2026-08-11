package user

import (
	"context"
	"strings"
	"time"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id string) (*UserResponse, error) {

	userId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	u, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*UserResponse, error) {

	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, ErrorEmailRequired
	}

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (*UserResponse, error) {

	phone = strings.TrimSpace(phone)
	if phone == "" {
		return nil, ErrorPhoneRequired
	}

	u, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func toUserResponse(u *domain.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:              u.ID.String(),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		EmailVerified:   u.EmailVerified,
		Phone:           u.Phone,
		PhoneVerified:   u.PhoneVerified,
		ProfileImageURL: u.ProfileImageURL,
		Role:            string(u.Role),
		Status:          string(u.Status),
		LastLoginAt:     u.LastLoginAt.Format(time.RFC3339),
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
	}
}
