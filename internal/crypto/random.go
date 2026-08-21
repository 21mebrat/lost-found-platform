package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateOTPCode() (string, error) {
	max := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate random number: %w", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func GenerateRandomToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random bytes for token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
