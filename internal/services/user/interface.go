package user

import "context"

type UserService interface {
	Register(
		ctx context.Context,
		req RegisterRequest,
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

	Update(
		ctx context.Context,
		id string,
		req UpdateRequest,
	) (*UserResponse, error)

	// Delete(
	// 	ctx context.Context,
	// 	id string,
	// ) error
}
