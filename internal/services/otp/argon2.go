package otp

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Memory      uint32 = 64 * 1024
	argon2Iterations  uint32 = 3
	argon2Parallelism uint8  = 2
	argon2SaltLength  uint32 = 16
	argon2KeyLength   uint32 = 32
)

type Argon2Hasher struct{}

var _ Hasher = (*Argon2Hasher)(nil)

func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{}
}

func (h *Argon2Hasher) Hash(
	ctx context.Context,
	code string,
) (string, error) {

	if code == "" {
		return "", errors.New("otp code cannot be empty")
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	salt := make([]byte, argon2SaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate otp salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(code),
		salt,
		argon2Iterations,
		argon2Memory,
		argon2Parallelism,
		argon2KeyLength,
	)

	return encodeHash(salt, hash), nil
}

func (h *Argon2Hasher) Compare(
	ctx context.Context,
	encodedHash string,
	code string,
) (bool, error) {

	if code == "" {
		return false, errors.New("otp code cannot be empty")
	}

	if err := ctx.Err(); err != nil {
		return false, err
	}

	salt, expectedHash, err := decodeHash(encodedHash)
	if err != nil {
		return false, fmt.Errorf("decode otp hash: %w", err)
	}

	actualHash := argon2.IDKey(
		[]byte(code),
		salt,
		argon2Iterations,
		argon2Memory,
		argon2Parallelism,
		argon2KeyLength,
	)

	if subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) != 1 {
		return false, nil
	}

	return true, nil
}

func encodeHash(salt, hash []byte) string {
	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory,
		argon2Iterations,
		argon2Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func decodeHash(encoded string) ([]byte, []byte, error) {
	var (
		version    int
		memory     uint32
		iterations uint32
		parallel   uint8
		saltBase64 string
		hashBase64 string
	)

	_, err := fmt.Sscanf(
		encoded,
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		&version,
		&memory,
		&iterations,
		&parallel,
		&saltBase64,
		&hashBase64,
	)

	if err != nil {
		return nil, nil, errors.New("invalid argon2 hash format")
	}

	if version != 19 {
		return nil, nil, errors.New("unsupported argon2 version")
	}

	if memory != argon2Memory ||
		iterations != argon2Iterations ||
		parallel != argon2Parallelism {
		return nil, nil, errors.New("unsupported argon2 parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, nil, errors.New("invalid argon2 salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(hashBase64)
	if err != nil {
		return nil, nil, errors.New("invalid argon2 hash")
	}

	if len(salt) != int(argon2SaltLength) {
		return nil, nil, errors.New("invalid argon2 salt length")
	}

	if len(hash) != int(argon2KeyLength) {
		return nil, nil, errors.New("invalid argon2 hash length")
	}

	return salt, hash, nil
}
