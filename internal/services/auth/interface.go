package auth

import "context"

type AuthService interface {
	Login(ctx context.Context, input LoginInput) (*LoginResponse, error)
	Refresh(ctx context.Context, accesstoken string) (*LoginResponse, error)
}
