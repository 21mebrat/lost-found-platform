package user

import (
	"context"
	"strings"
	"time"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrorUserNotFound
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

	if u == nil {
		return nil, ErrorUserNotFound
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

	if u == nil {
		return nil, ErrorUserNotFound
	}

	return toUserResponse(u), nil
}

func toUserResponse(user *domain.User) *UserResponse {
	response := &UserResponse{
		ID:              user.ID.String(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Phone:           user.Phone,
		EmailVerified:   user.EmailVerified,
		PhoneVerified:   user.PhoneVerified,
		Role:            string(user.Role),
		Status:          string(user.Status),
		ProfileImageURL: user.ProfileImageURL,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}

	if user.LastLoginAt != nil {
		response.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}

	return response
}
