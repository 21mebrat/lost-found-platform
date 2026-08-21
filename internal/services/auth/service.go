package auth

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/21mebrat/lost-found-platform/internal/auth"
	"github.com/21mebrat/lost-found-platform/internal/domain/token"
	tokenrepo "github.com/21mebrat/lost-found-platform/internal/repository/token"
	userrep "github.com/21mebrat/lost-found-platform/internal/repository/user"
	"github.com/21mebrat/lost-found-platform/internal/validator"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	token tokenrepo.Repository
	JWT   *auth.JWTManager
	users userrep.Repository
	redi  *redis.Client
}

// compile time check

var _ AuthService = (*Service)(nil)

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {

	password := strings.TrimSpace(input.Password)
	email := strings.TrimSpace(input.Email)

	if password == "" || email == "" {
		return nil, ErrorIncorrectCredential
	}
	// validate email
	_, err := mail.ParseAddress(input.Email)

	if err != nil {
		return nil, ErrorIncorrectCredential
	}
	user, err := s.users.GetByEmail(ctx, email)

	if err != nil || user == nil {
		return nil, ErrorIncorrectCredential
	}

	if !validator.ValidatePassword(input.Password) {
		return nil, ErrorIncorrectCredential
	}

	err = auth.ComparePassword(user.PasswordHash, strings.TrimSpace(input.Password))
	if err != nil {
		return nil, err
	}
	// access token
	accessToken, expiresAt, err := s.JWT.GenerateAccessToken(user.ID, *user.Email, string(user.Role))

	if err != nil {
		return nil, err
	}

	// refrash token
	refreshToken, refreshExpiresAt, err := s.JWT.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// store refrash token
	refresh := token.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExpiresAt,
	}
	// create access token
	err = s.token.Create(ctx, &refresh)
	if err != nil {
		return nil, err
	}

	expireIn := int16(
		time.Until(expiresAt).Seconds(),
	)

	return &LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresIn:    uint64(expireIn),
	}, nil

}

func (s *Service) Refresh(ctx context.Context, token string) (*LoginResponse, error) {
	return nil, nil
}
