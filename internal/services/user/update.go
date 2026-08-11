package user

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (*UserResponse, error) {

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrorUserNotFound
	}

	existingUser, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrorUserNotFound
	}

	if input.FirstName != "" {
		firstName := strings.TrimSpace(input.FirstName)
		if len(firstName) < 2 {
			return nil, ErrorFirstNameRequired
		}
		existingUser.FirstName = firstName
	}

	if input.LastName != "" {
		lastName := strings.TrimSpace(input.LastName)
		if len(lastName) < 2 {
			return nil, ErrorLastNameRequired
		}
		existingUser.LastName = lastName
	}

	if input.ProfileImageURL != "" {
		existingUser.ProfileImageURL = strings.TrimSpace(input.ProfileImageURL)
	}

	existingUser.UpdatedAt = time.Now().UTC()

	updatedUser, err := s.repo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:              updatedUser.ID.String(),
		FirstName:       updatedUser.FirstName,
		LastName:        updatedUser.LastName,
		Email:           updatedUser.Email,
		Phone:           updatedUser.Phone,
		EmailVerified:   updatedUser.EmailVerified,
		PhoneVerified:   updatedUser.PhoneVerified,
		ProfileImageURL: updatedUser.ProfileImageURL,
		Role:            string(updatedUser.Role),
		Status:          string(updatedUser.Status),
	}, nil
}
