package otp

import (
	"context"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/otp"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session *domain.OTPSession) error
	Get(ctx context.Context, id uuid.UUID) (*domain.OTPSession, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) (int, error)
}
