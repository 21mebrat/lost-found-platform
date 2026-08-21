package user

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/validator"
	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (*UserResponse, error) {
	userID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, ErrorInvalidID
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

	if input.MiddleName != "" {
		middleName := strings.TrimSpace(input.MiddleName)
		if len(middleName) < 2 {
			return nil, ErrorMiddleNameRequired
		}
		existingUser.MiddleName = middleName
	}

	if input.LastName != "" {
		lastName := strings.TrimSpace(input.LastName)
		if len(lastName) < 2 {
			return nil, ErrorLastNameRequired
		}
		existingUser.LastName = lastName
	}

	if input.Phone != "" {
		phone, err := validator.ValidatePhone(input.Phone)
		if err != nil {
			return nil, err
		}
		if phone != existingUser.Phone {
			existingByPhone, err := s.repo.GetByPhone(ctx, phone)
			if err != nil && !errors.Is(err, ErrorUserNotFound) {
				return nil, err
			}
			if existingByPhone != nil {
				return nil, ErrorPhoneAlreadyExists
			}
			existingUser.Phone = phone
			existingUser.PhoneVerified = false
		}
	}

	if input.Email != nil && strings.TrimSpace(*input.Email) != "" {
		cleanEmail := strings.ToLower(strings.TrimSpace(*input.Email))
		_, err := mail.ParseAddress(cleanEmail)
		if err != nil {
			return nil, ErrorInvalidEmail
		}
		if existingUser.Email == nil || *existingUser.Email != cleanEmail {
			existingByEmail, err := s.repo.GetByEmail(ctx, cleanEmail)
			if err != nil && !errors.Is(err, ErrorUserNotFound) {
				return nil, err
			}
			if existingByEmail != nil {
				return nil, ErrorEmailAlreadyExists
			}
			existingUser.Email = &cleanEmail
		}
	}

	if input.Fayda != nil && strings.TrimSpace(*input.Fayda) != "" {
		cleanFayda := strings.TrimSpace(*input.Fayda)
		if existingUser.Fayda == nil || *existingUser.Fayda != cleanFayda {
			existingByFayda, err := s.repo.GetByFayda(ctx, cleanFayda)
			if err != nil && !errors.Is(err, ErrorUserNotFound) {
				return nil, err
			}
			if existingByFayda != nil {
				return nil, ErrorFaydaAlreadyExists
			}
			existingUser.Fayda = &cleanFayda
			existingUser.FaydaVerified = false
		}
	}

	if input.LanguageCode != "" {
		lang := strings.ToLower(strings.TrimSpace(input.LanguageCode))
		existingUser.LanguageCode = lang
	}

	if input.ProfileImageURL != nil {
		imgURL := strings.TrimSpace(*input.ProfileImageURL)
		existingUser.ProfileImageURL = &imgURL
	}

	existingUser.UpdatedAt = time.Now().UTC()

	updatedUser, err := s.repo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	return ToUserResponse(updatedUser), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return ErrorInvalidID
	}

	return s.repo.Delete(ctx, userID)
}
