package user

import (
	"context"

	userrepo "github.com/21mebrat/lost-found-platform/internal/repository/user"
)

type UserService interface {
	SendOTP(
		ctx context.Context,
		req SendOTPRequest,
	) (*SendOTPResponse, error)

	VerifyOTP(
		ctx context.Context,
		req VerifyOTPRequest,
	) (*VerifyOTPResponse, error)

	CompleteProfile(
		ctx context.Context,
		req CompleteProfileRequest,
	) (*UserResponse, error)

	GetByID(
		ctx context.Context,
		id string,
	) (*UserResponse, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*UserResponse, error)

	GetByPhone(
		ctx context.Context,
		phone string,
	) (*UserResponse, error)

	GetByFayda(
		ctx context.Context,
		fayda string,
	) (*UserResponse, error)

	GetList(
		ctx context.Context,
		params userrepo.ListParams,
	) (*UserListResponse, error)

	Update(
		ctx context.Context,
		id string,
		req UpdateRequest,
	) (*UserResponse, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}
