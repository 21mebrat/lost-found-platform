package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`

	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`

	jwt.RegisteredClaims
}
