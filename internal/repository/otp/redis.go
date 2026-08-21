package otp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	domain "github.com/21mebrat/lost-found-platform/internal/domain/otp"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	otpKeyPrefix = "otp:session:"
)

type RedisRepository struct {
	client *redis.Client
}

var _ Repository = (*RedisRepository)(nil)

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func otpKey(id uuid.UUID) string {
	return otpKeyPrefix + id.String()
}

func (r *RedisRepository) Create(
	ctx context.Context,
	session *domain.OTPSession,
) error {

	if session == nil {
		return errors.New("otp session cannot be nil")
	}

	if session.ID == uuid.Nil {
		return errors.New("otp session id cannot be nil")
	}

	ttl := time.Until(session.ExpiresAt)

	if ttl <= 0 {
		return errors.New("otp session already expired")
	}

	key := otpKey(session.ID)

	fields := map[string]any{
		"phone":      session.Phone,
		"code_hash":  session.CodeHash,
		"purpose":    string(session.Purpose),
		"attempts":   session.Attempts,
		"created_at": session.CreatedAt.UTC().Format(time.RFC3339Nano),
		"expires_at": session.ExpiresAt.UTC().Format(time.RFC3339Nano),
	}

	if err := r.client.HSet(
		ctx,
		key,
		fields,
	).Err(); err != nil {
		return fmt.Errorf("set otp session: %w", err)
	}

	if err := r.client.Expire(
		ctx,
		key,
		ttl,
	).Err(); err != nil {
		// The hash was created, but without TTL it could remain forever.
		// Remove it so we don't leave an unsafe/stale OTP session.
		_ = r.client.Del(ctx, key).Err()

		return fmt.Errorf("set otp session ttl: %w", err)
	}

	return nil
}

func (r *RedisRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*domain.OTPSession, error) {

	key := otpKey(id)

	fields, err := r.client.HGetAll(
		ctx,
		key,
	).Result()

	if err != nil {
		return nil, fmt.Errorf("get otp session: %w", err)
	}

	if len(fields) == 0 {
		return nil, nil
	}

	session, err := parseSession(id, fields)
	if err != nil {
		return nil, fmt.Errorf("parse otp session: %w", err)
	}

	return session, nil
}

func (r *RedisRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	key := otpKey(id)

	if err := r.client.Del(
		ctx,
		key,
	).Err(); err != nil {
		return fmt.Errorf("delete otp session: %w", err)
	}

	return nil
}

func (r *RedisRepository) IncrementAttempts(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {

	key := otpKey(id)

	attempts, err := r.client.HIncrBy(
		ctx,
		key,
		"attempts",
		1,
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, domain.ErrOTPSessionNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("increment otp attempts: %w", err)
	}

	return int(attempts), nil
}

func parseSession(
	id uuid.UUID,
	fields map[string]string,
) (*domain.OTPSession, error) {

	attempts, err := strconv.Atoi(fields["attempts"])
	if err != nil {
		return nil, fmt.Errorf("parse attempts: %w", err)
	}

	createdAt, err := time.Parse(
		time.RFC3339Nano,
		fields["created_at"],
	)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}

	expiresAt, err := time.Parse(
		time.RFC3339Nano,
		fields["expires_at"],
	)
	if err != nil {
		return nil, fmt.Errorf("parse expires_at: %w", err)
	}

	return &domain.OTPSession{
		ID:        id,
		Phone:     fields["phone"],
		CodeHash:  fields["code_hash"],
		Purpose:   domain.Purpose(fields["purpose"]),
		Attempts:  attempts,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}
