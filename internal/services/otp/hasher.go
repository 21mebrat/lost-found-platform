package otp

import "context"

type Hasher interface {
	Hash(
		ctx context.Context,
		code string,
	) (string, error)

	Compare(
		ctx context.Context,
		hash string,
		code string,
	) (bool, error)
}
