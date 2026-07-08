package services

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode generates a random 6-character short code.
func GenerateShortCode() (string, error) {

	length := 6
	shortCode := make([]byte, length)

	for i := range shortCode {

		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		shortCode[i] = charset[randomIndex.Int64()]
	}

	return string(shortCode), nil
}
