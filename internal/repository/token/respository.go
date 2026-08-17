package token

import (
	"context"

	"github.com/21mebrat/lost-found-platform/internal/domain/token"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, token *token.RefreshToken) error
	Get(ctx context.Context, token string) (*token.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	Delete(ctx context.Context, token string) error
	DeleteByUserId(ctx context.Context, user_id uuid.UUID) error
}
