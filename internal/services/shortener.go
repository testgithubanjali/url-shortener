package services

import (
	"crypto/rand"
	"log"
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
			log.Printf("Failed to generate random bytes for short code: %v", err)
			return "", err
		}

		shortCode[i] = charset[randomIndex.Int64()]
	}

	generatedCode := string(shortCode)
	log.Printf("Generated short code %s", generatedCode)
	return generatedCode, nil
}
