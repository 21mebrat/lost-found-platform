package user

import (
	"context"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/user"
	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		user *domain.User,
	) (*domain.User, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.User, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*domain.User, error)

	GetByFayda(
		ctx context.Context,
		fayda string,
	) (*domain.User, error)

	GetByPhone(
		ctx context.Context,
		phone string,
	) (*domain.User, error)

	Update(
		ctx context.Context,
		user *domain.User,
	) (*domain.User, error)

	GetList(
		ctx context.Context,
		params ListParams,
	) (*PaginatedUsers, error)

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
