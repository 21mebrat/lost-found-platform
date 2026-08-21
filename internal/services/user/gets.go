package user

import (
	"context"
	"strings"
	"time"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	userID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, ErrorInvalidID
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrorUserNotFound
	}

	return ToUserResponse(u), nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*UserResponse, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	if cleanEmail == "" {
		return nil, ErrorEmailRequired
	}

	u, err := s.repo.GetByEmail(ctx, cleanEmail)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrorUserNotFound
	}

	return ToUserResponse(u), nil
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (*UserResponse, error) {
	cleanPhone := strings.TrimSpace(phone)

	if cleanPhone == "" {
		return nil, ErrorPhoneRequired
	}

	u, err := s.repo.GetByPhone(ctx, cleanPhone)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrorUserNotFound
	}

	return ToUserResponse(u), nil
}

func (s *Service) GetByFayda(ctx context.Context, fayda string) (*UserResponse, error) {
	cleanFayda := strings.TrimSpace(fayda)

	if cleanFayda == "" {
		return nil, ErrorUserNotFound
	}

	u, err := s.repo.GetByFayda(ctx, cleanFayda)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrorUserNotFound
	}

	return ToUserResponse(u), nil
}

func (s *Service) GetList(ctx context.Context, params userrepo.ListParams) (*UserListResponse, error) {
	paginated, err := s.repo.GetList(ctx, params)
	if err != nil {
		return nil, err
	}

	return ToUserListResponse(paginated), nil
}

func ToUserResponse(u *domain.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:              u.ID.String(),
		FirstName:       u.FirstName,
		MiddleName:      u.MiddleName,
		LastName:        u.LastName,
		Phone:           u.Phone,
		Email:           u.Email,
		Fayda:           u.Fayda,
		LanguageCode:    u.LanguageCode,
		PhoneVerified:   u.PhoneVerified,
		FaydaVerified:   u.FaydaVerified,
		ProfileImageURL: u.ProfileImageURL,
		Status:          string(u.Status),
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       u.UpdatedAt.Format(time.RFC3339),
	}
}
