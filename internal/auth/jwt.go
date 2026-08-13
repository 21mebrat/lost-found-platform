package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secret               string
	issuer               string
	refreshTokenDuration time.Duration
	accessTokenDuration  time.Duration
}

func NewJWTManager(
	secret, issuer string,
	refreshTokenDuration,
	accessTokenDuration time.Duration,
) *JWTManager {
	return &JWTManager{
		secret:               secret,
		issuer:               issuer,
		refreshTokenDuration: refreshTokenDuration,
		accessTokenDuration:  accessTokenDuration,
	}
}

func (j *JWTManager) GenerateAccessToken(id uuid.UUID, email, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.accessTokenDuration)
	claim := &AccessTokenClaims{
		UserID: id,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  j.issuer,
			Subject: id.String(),
			Audience: []string{
				"web app",
				"mobile app",
			},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim,
	)
	signedSignature, err := token.SignedString([]byte(j.secret))

	if err != nil {
		return "", time.Time{}, err
	}

	return signedSignature, expiresAt, nil
}

// verify token
func (j *JWTManager) VerifyToken(
	tokenString string,
) (*AccessTokenClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaims{},
		func(t *jwt.Token) (any, error) {

			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return []byte(j.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// generate refresh token
func (j *JWTManager) GenerateRefreshToken(id uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.refreshTokenDuration)

	claims := &RefreshTokenClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  j.issuer,
			Subject: id.String(),
			Audience: []string{
				"web app",
				"mobile app",
			},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(j.secret))

	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresAt, nil
}
